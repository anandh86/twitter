package handlers

type UserResponseWithTokenDTO struct {
	Email        string `json:"email"`
	ID           int    `json:"id"`
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type UserResponseDTO struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
}

type UserRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
