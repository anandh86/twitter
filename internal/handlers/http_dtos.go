package handlers

type UserResponseWithTokenDTO struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
	Token string `json:"token"`
}

type UserResponseDTO struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
}

type UserRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Expiry   int    `json:"expires_in_seconds"`
}
