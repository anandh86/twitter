package ports

import "github.com/anandhmaps/chirpy/internal/core/domain"

// IUseCase is a primary port that the core must respond to
type IUseCase interface {
	CreateUser(emailid string, password string) (domain.User, error)
	UpdateUser(id int, emailid string, password string) error
	GetUserById(id int) (domain.User, error)
	LoginUser(emailid string, password string) (int, error)
	PostTweet(body string) (domain.Tweet, error)
	GetTweetById(id int) (domain.Tweet, error)
}
