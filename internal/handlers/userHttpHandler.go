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

	userRequest := domain.User{}

	err := decoder.Decode(&userRequest)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	userResponse, _ := u.uuc.CreateUser(userRequest.Email)
	respondWithJSON(w, http.StatusCreated, userResponse)
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

	tweetResponse, _ := u.uuc.GetTweetById(chirpID)
	respondWithJSON(w, http.StatusCreated, tweetResponse)
}
