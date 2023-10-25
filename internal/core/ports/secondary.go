package ports

import "github.com/anandhmaps/chirpy/internal/core/domain"

// UserRepository is a secondary port that the core will make calls to
type UserRepository interface {
	Save(user domain.User) (domain.User, error)
	GetUserById(id int) (domain.User, error)
}
