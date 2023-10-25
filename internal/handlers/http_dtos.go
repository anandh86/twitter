package handlers

type UserResponseDTO struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
}

type UserRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
