package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/anandhmaps/chirpy/internal/core/domain"
	"github.com/anandhmaps/chirpy/internal/core/ports"
	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func ProvideUserHttpHandler(uuc ports.IUseCase) *UserHttpHandler {
	// by default, godotenv will look for a file named .env in the current directory
	godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")

	return &UserHttpHandler{
		uuc:   uuc,
		token: jwtSecret,
	}
}

type UserHttpHandler struct {
	uuc   ports.IUseCase
	token string
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
	expiresAt := newFunction(expiresInSeconds)

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

// adjust expiry time based on biz logic
func newFunction(expiresInSeconds int) time.Time {
	var expiresAt time.Time
	defaultExpirationTime := 24 * time.Hour

	if expiresInSeconds == 0 {
		expiresAt = time.Now().Add(defaultExpirationTime)
	} else if expiresInSeconds > int(defaultExpirationTime.Seconds()) {
		expiresAt = time.Now().Add(defaultExpirationTime)
	} else {
		expiresAt = time.Now().Add(time.Duration(expiresInSeconds) * time.Second)
	}
	return expiresAt
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
	oneHourInSec := 60 * 60
	accessToken, _ := u.createJWTToken(userId, oneHourInSec, "chirpy-access")

	// create refresh token
	sixtyDaysInSec := 60 * 60 * 24 * 60
	refreshToken, _ := u.createJWTToken(userId, sixtyDaysInSec, "chirpy-refresh")

	// presentation segment
	userResponseDTO := UserResponseWithTokenDTO{
		ID:           userId,
		Email:        userRequest.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	respondWithJSON(w, http.StatusOK, userResponseDTO)
}

func fetchBearerToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")

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

	decoder := json.NewDecoder(r.Body)
	tweetRequest := domain.Tweet{}

	err := decoder.Decode(&tweetRequest)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	tweetResponse, _ := u.uuc.PostTweet(tweetRequest.Body)
	respondWithJSON(w, http.StatusCreated, tweetResponse)
}

func (u *UserHttpHandler) GetTweetById(w http.ResponseWriter, r *http.Request) {

	// Read the input parameter
	chirpIDStr := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(chirpIDStr)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid parameters")
		return
	}

	tweetResponse, err1 := u.uuc.GetTweetById(chirpID)

	if err1 != nil {
		respondWithError(w, http.StatusNotFound, "Invalid parameters")
	}

	respondWithJSON(w, http.StatusOK, tweetResponse)
}
