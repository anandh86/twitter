package adapters

import (
	"errors"

	"github.com/anandhmaps/chirpy/internal/core/domain"
	"github.com/anandhmaps/chirpy/internal/core/ports"
)

// In memory implementation
func ProvideUserRepository() ports.UserRepository {
	return &myInMemoryUserRepository{
		userMap:          make(map[int]domain.User),
		currentNoOfUsers: 0,
	}
}

// myInMemoryUserRepository implements ports.UserRepository
type myInMemoryUserRepository struct {
	userMap          map[int]domain.User
	currentNoOfUsers int
}

func (u *myInMemoryUserRepository) Save(user domain.User) (domain.User, error) {

	if _, ok := u.userMap[user.ID]; ok {
		// user already present
		return user, errors.ErrUnsupported
	}

	newID := u.currentNoOfUsers + 1
	u.currentNoOfUsers = newID
	user.ID = newID

	u.userMap[newID] = user

	return user, nil
}

func (u *myInMemoryUserRepository) GetUserById(id int) (domain.User, error) {
	user, ok := u.userMap[id]

	if !ok {
		return user, errors.New("user id not found")
	}

	return user, nil
}
