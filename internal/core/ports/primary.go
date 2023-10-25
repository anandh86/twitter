package ports

import "github.com/anandhmaps/chirpy/internal/core/domain"

// IUseCase is a primary port that the core must respond to
type IUseCase interface {
	CreateUser(emailid string) (domain.User, error)
	GetUserById(id int) (domain.User, error)
	PostTweet(body string) (domain.Tweet, error)
	GetTweetById(id int) (domain.Tweet, error)
}
