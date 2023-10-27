package handlers

type UserResponseWithTokenDTO struct {
	Email        string `json:"email"`
	ID           int    `json:"id"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type UserResponseDTO struct {
	Email       string `json:"email"`
	ID          int    `json:"id"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type UserRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Data struct {
	UserID int `json:"user_id"`
}

type WebHookBody struct {
	Event string `json:"event"`
	Data  Data   `json:"data"`
}
