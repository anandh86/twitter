package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/anandhmaps/chirpy/internal/core/domain"
	"github.com/anandhmaps/chirpy/internal/core/ports"
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
