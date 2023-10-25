package usecases

import (
	"github.com/anandhmaps/chirpy/internal/core/domain"
	"github.com/anandhmaps/chirpy/internal/core/ports"
)

func ProvideUserUseCase(userRepository ports.UserRepository) ports.UserUseCase {
	return &userUseCase{
		userRepository: userRepository,
	}
}

// userUseCase implements ports.UserUseCase
type userUseCase struct {
	userRepository ports.UserRepository
}

func (u userUseCase) CreateUser(emailid string) (domain.User, error) {
	user := domain.User{Email: emailid}
	return u.userRepository.Save(user)
}

func (u userUseCase) GetUserById(id int) (domain.User, error) {
	return u.userRepository.GetUserById(id)
}
