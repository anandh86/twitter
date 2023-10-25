package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/anandhmaps/chirpy/internal/core/domain"
	"github.com/anandhmaps/chirpy/internal/core/ports"
	"github.com/go-chi/chi"
)

func ProvideUserHttpHandler(uuc ports.IUseCase) *UserHttpHandler {
	return &UserHttpHandler{
		uuc: uuc,
	}
}

type UserHttpHandler struct {
	uuc ports.IUseCase
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
	}

	// success
	userResponseDTO := UserResponseDTO{
		ID:    userId,
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
