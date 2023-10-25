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

func (u *UserHttpHandler) createJWTToken(userId int, expiresInSeconds int) (string, error) {

	mySigningKey := []byte(u.token)
	var expiresAt time.Time

	// expiry time calculation
	{
		defaultExpirationTime := 24 * time.Hour

		if expiresInSeconds == 0 {
			expiresAt = time.Now().Add(defaultExpirationTime)
		} else if expiresInSeconds > int(defaultExpirationTime.Seconds()) {
			expiresAt = time.Now().Add(defaultExpirationTime)
		} else {
			expiresAt = time.Now().Add(time.Duration(expiresInSeconds) * time.Second)
		}
	}

	// create the claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   strconv.Itoa(userId),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

func (u *UserHttpHandler) LoginUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	userRequest := UserRequestDTO{}

	err := decoder.Decode(&userRequest)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	userId, err1 := u.uuc.LoginUser(userRequest.Email, userRequest.Password)

	if err1 != nil {
		// Login failed
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// success
	token, _ := u.createJWTToken(userId, userRequest.Expiry)

	userResponseDTO := UserResponseWithTokenDTO{
		ID:    userId,
		Email: userRequest.Email,
		Token: token,
	}
	respondWithJSON(w, http.StatusOK, userResponseDTO)
}

func (u *UserHttpHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	bearerToken := r.Header.Get("Authorization")

	parts := strings.Fields(bearerToken)

	tokenString := parts[1]

	{
		tokenSecret := u.token
		claimsStruct := jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			&claimsStruct,
			func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
		)

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		if !token.Valid {
			respondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		// if valid token, update the user details
		userId, _ := token.Claims.GetSubject()
		userIdInt, _ := strconv.Atoi(userId)

		decoder := json.NewDecoder(r.Body)

		userRequest := UserRequestDTO{}

		errDecode := decoder.Decode(&userRequest)

		if errDecode != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
			return
		}

		if authErr := u.uuc.UpdateUser(userIdInt, userRequest.Email, userRequest.Password); authErr != nil {
			respondWithError(w, http.StatusUnauthorized, "got issues")
		}

		userResponseDTO := UserResponseDTO{
			ID:    userIdInt,
			Email: userRequest.Email,
		}
		respondWithJSON(w, http.StatusOK, userResponseDTO)
	}
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
