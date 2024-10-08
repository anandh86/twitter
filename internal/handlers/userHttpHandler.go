package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/anandh86/chirpy/internal/core/domain"
	"github.com/anandh86/chirpy/internal/core/ports"
	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func ProvideUserHttpHandler(uuc ports.IUseCase) *UserHttpHandler {
	// by default, godotenv will look for a file named .env in the current directory
	godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")
	apiKey := os.Getenv("POLKA_KEY")

	return &UserHttpHandler{
		uuc:         uuc,
		token:       jwtSecret,
		polkaApiKey: apiKey,
	}
}

type UserHttpHandler struct {
	uuc         ports.IUseCase
	token       string
	polkaApiKey string
}

func (u *UserHttpHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	userRequest := UserRequestDTO{}

	err := decoder.Decode(&userRequest)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	userResponse, errCreation := u.uuc.CreateUser(userRequest.Email, userRequest.Password)

	if errCreation != nil {
		respondWithError(w, http.StatusBadRequest, "account present already")
		return
	}
	userResponseDTO := UserResponseDTO{
		ID:    userResponse.ID,
		Email: userResponse.Email,
	}
	respondWithJSON(w, http.StatusCreated, userResponseDTO)
}

func (u *UserHttpHandler) createJWTToken(userId int, expiresInSeconds int, issuer string) (string, error) {
	mySigningKey := []byte(u.token)
	expiresAt := time.Now().Add(time.Duration(expiresInSeconds) * time.Second)

	// create the claims
	claims := &jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   strconv.Itoa(userId),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

func (u *UserHttpHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	tokenString := fetchBearerToken(r)

	if !u.uuc.RevokeRefreshToken(tokenString) {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	respondWithJSON(w, http.StatusOK, "")

}

func (u *UserHttpHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	tokenString := fetchBearerToken(r)

	isValidToken, jwtToken := isValidToken(tokenString, u.token)

	if !isValidToken {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	tokenIssuer, _ := jwtToken.Claims.GetIssuer()

	if tokenIssuer != "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if u.uuc.IsRefreshTokenRevoked(tokenString) {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// generate new access token
	userIdStr, _ := jwtToken.Claims.GetSubject()
	userId, _ := strconv.Atoi(userIdStr)

	accessToken, _ := u.generateAccessToken(userId)

	response := map[string]string{
		"token": accessToken,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (u *UserHttpHandler) generateAccessToken(userId int) (string, error) {
	issuer := "chirpy-access"
	// access tokens should expire in one hour
	expiryInSeconds := int(time.Hour.Seconds())

	return u.createJWTToken(userId, expiryInSeconds, issuer)
}

func (u *UserHttpHandler) generateRefreshToken(userId int) (string, error) {
	issuer := "chirpy-refresh"
	// refresh tokens should expire in 60 days
	expireIn60Days := 60 * 24 * time.Hour

	return u.createJWTToken(userId, int(expireIn60Days.Seconds()), issuer)
}

func (u *UserHttpHandler) LoginUser(w http.ResponseWriter, r *http.Request) {

	// fetch input info
	decoder := json.NewDecoder(r.Body)
	userRequest := UserRequestDTO{}
	err := decoder.Decode(&userRequest)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed json body")
		return
	}

	// call use case
	userId, err1 := u.uuc.LoginUser(userRequest.Email, userRequest.Password)

	if err1 != nil {
		// Login failed
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// create access token
	accessToken, _ := u.generateAccessToken(userId)

	// create refresh token
	refreshToken, _ := u.generateRefreshToken(userId)
	u.uuc.StoreRefreshToken(refreshToken)

	repoUser, errFunc := u.uuc.GetUserById(userId)

	if errFunc != nil {
		respondWithError(w, http.StatusBadRequest, "Server issues")
		return
	}

	// presentation segment
	userResponseDTO := UserResponseWithTokenDTO{
		ID:           userId,
		Email:        userRequest.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IsChirpyRed:  repoUser.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusOK, userResponseDTO)
}

func fetchApiKey(r *http.Request) string {
	authorizationHeaderToken := r.Header.Get("Authorization")

	if authorizationHeaderToken == "" {
		return ""
	}

	parts := strings.Fields(authorizationHeaderToken)

	if parts[0] != "ApiKey" {
		return ""
	}

	return parts[1]
}

func fetchBearerToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")

	if bearerToken == "" {
		return bearerToken
	}

	parts := strings.Fields(bearerToken)

	return parts[1]
}

func isValidToken(jwtToken string, tokenSecret string) (bool, *jwt.Token) {

	claimsStruct := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		jwtToken,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)

	return (err == nil && token.Valid), token
}

func (u *UserHttpHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	tokenString := fetchBearerToken(r)

	isValidToken, jwtToken := isValidToken(tokenString, u.token)

	if !isValidToken {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	tokenIssuer, _ := jwtToken.Claims.GetIssuer()

	if tokenIssuer != "chirpy-access" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// update the user info
	userId, _ := jwtToken.Claims.GetSubject()
	userIdInt, _ := strconv.Atoi(userId)

	decoder := json.NewDecoder(r.Body)
	userRequest := UserRequestDTO{}
	errDecode := decoder.Decode(&userRequest)

	if errDecode != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed json body")
		return
	}

	if authErr := u.uuc.UpdateUser(userIdInt, userRequest.Email, userRequest.Password); authErr != nil {
		respondWithError(w, http.StatusUnauthorized, "got issues")
		return
	}

	// presentation segment
	userResponseDTO := UserResponseDTO{
		ID:    userIdInt,
		Email: userRequest.Email,
	}
	respondWithJSON(w, http.StatusOK, userResponseDTO)
}

func (u *UserHttpHandler) PostTweet(w http.ResponseWriter, r *http.Request) {

	// authenticate the user first
	tokenString := fetchBearerToken(r)

	isValidToken, jwtToken := isValidToken(tokenString, u.token)

	if !isValidToken {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	tokenIssuer, _ := jwtToken.Claims.GetIssuer()

	if tokenIssuer != "chirpy-access" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// figure out author's id
	authorIdStr, _ := jwtToken.Claims.GetSubject()
	authorId, _ := strconv.Atoi(authorIdStr)

	decoder := json.NewDecoder(r.Body)
	tweetRequest := domain.Tweet{}

	err := decoder.Decode(&tweetRequest)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	tweetResponse, _ := u.uuc.PostTweet(tweetRequest.Body, authorId)
	tweetResponse.AuthorId = authorId
	respondWithJSON(w, http.StatusCreated, tweetResponse)
}

func (u *UserHttpHandler) DeleteTweet(w http.ResponseWriter, r *http.Request) {
	// Read the input parameter
	tweetIdStr := chi.URLParam(r, "tweetId")
	tweetId, err := strconv.Atoi(tweetIdStr)

	// authenticate the user first
	tokenString := fetchBearerToken(r)

	isValidToken, jwtToken := isValidToken(tokenString, u.token)

	if !isValidToken {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	tokenIssuer, _ := jwtToken.Claims.GetIssuer()

	if tokenIssuer != "chirpy-access" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// figure out author's id
	authorIdStr, _ := jwtToken.Claims.GetSubject()
	authorId, _ := strconv.Atoi(authorIdStr)

	if err != nil || tweetId == 0 {
		respondWithError(w, http.StatusInternalServerError, "Invalid parameters")
		return
	}

	if deletionErr := u.uuc.DeleteTweet(tweetId, authorId); deletionErr != nil {
		respondWithError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	respondWithJSON(w, http.StatusOK, "all good")
}

func (u *UserHttpHandler) GetTweetById(w http.ResponseWriter, r *http.Request) {

	// Read the input parameter
	tweetIdStr := chi.URLParam(r, "tweetId")
	tweetId, err := strconv.Atoi(tweetIdStr)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid parameters")
		return
	}

	tweetResponse, err1 := u.uuc.GetTweetById(tweetId)

	if err1 != nil {
		respondWithError(w, http.StatusNotFound, "Invalid parameters")
	}

	respondWithJSON(w, http.StatusOK, tweetResponse)
}

func (u *UserHttpHandler) GetAllTweets(w http.ResponseWriter, r *http.Request) {

	// s is a string that contains the value of the author_id query parameter
	// if it exists, or an empty string if it doesn't
	authorIdStr := r.URL.Query().Get("author_id")
	var allTweets []domain.Tweet

	if authorIdStr == "" {
		allTweets, _ = u.uuc.GetAllTweets()
	} else {
		author_id, _ := strconv.Atoi(authorIdStr)
		allTweets, _ = u.uuc.GetAuthorTweets(author_id)
	}

	sortBy := r.URL.Query().Get("sort")

	if sortBy == "desc" {
		sort.Sort(ByTweetIdDesc(allTweets))
	} else {
		sort.Sort(ByTweetIdAsc(allTweets))
	}

	respondWithJSON(w, http.StatusOK, allTweets)
}

func (u *UserHttpHandler) PolkaWebHooks(w http.ResponseWriter, r *http.Request) {

	apiKey := fetchApiKey(r)

	if apiKey != u.polkaApiKey {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	webHook := WebHookBody{}
	err := decoder.Decode(&webHook)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if webHook.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, "ok")
		return
	}

	if errMembership := u.uuc.UpdateUserMembership(webHook.Data.UserID, true); errMembership != nil {
		respondWithError(w, http.StatusNotFound, "user not found")
		return
	}

	respondWithJSON(w, http.StatusOK, "membership updated")

}

// Define a type for sorting by TweetId in ascending order
type ByTweetIdAsc []domain.Tweet

func (a ByTweetIdAsc) Len() int           { return len(a) }
func (a ByTweetIdAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTweetIdAsc) Less(i, j int) bool { return a[i].TweetId < a[j].TweetId }

// Define a type for sorting by TweetId in descending order
type ByTweetIdDesc []domain.Tweet

func (a ByTweetIdDesc) Len() int           { return len(a) }
func (a ByTweetIdDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTweetIdDesc) Less(i, j int) bool { return a[i].TweetId > a[j].TweetId }
