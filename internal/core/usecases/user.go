package usecases

import (
	"github.com/anandhmaps/chirpy/internal/core/domain"
	"github.com/anandhmaps/chirpy/internal/core/ports"
)

func ProvideUserUseCase(repoImplementation ports.IRepository) ports.IUseCase {
	return &userUseCase{
		repoImpl: repoImplementation,
	}
}

// userUseCase implements ports.UserUseCase
type userUseCase struct {
	repoImpl ports.IRepository
}

func (u userUseCase) CreateUser(emailid string) (domain.User, error) {
	user := domain.User{Email: emailid}
	return u.repoImpl.Save(user)
}

func (u userUseCase) GetUserById(id int) (domain.User, error) {
	return u.repoImpl.GetUserById(id)
}
