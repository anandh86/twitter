package ports

import "github.com/anandhmaps/chirpy/internal/core/domain"

// UserUseCase is a primary port that the core must respond to
type UserUseCase interface {
	CreateUser(emailid string) (domain.User, error)
	GetUserById(id int) (domain.User, error)
}
