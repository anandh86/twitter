package ports

import "github.com/anandhmaps/chirpy/internal/core/domain"

// IRepository is a secondary port that the core will make calls to
type IRepository interface {
	Save(user domain.User) (domain.User, error)
	GetUserById(id int) (domain.User, error)
	SaveTweet(tweet domain.Tweet) (domain.Tweet, error)
	GetTweetById(id int) (domain.Tweet, error)
}
