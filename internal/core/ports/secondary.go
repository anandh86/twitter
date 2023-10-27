package ports

import "github.com/anandhmaps/chirpy/internal/core/domain"

// IRepository is a secondary port that the core will make calls to
type IRepository interface {
	Save(user domain.User) (domain.User, error)
	GetUserById(id int) (domain.User, error)
	GetUserId(emailid string) (int, error)
	UpdateUser(id int, user domain.User) error
	UpdateUserMembership(id int, isMember bool) error
	SaveTweet(tweet domain.Tweet) (domain.Tweet, error)
	DeleteTweet(tweet domain.Tweet) error
	GetTweetById(id int) (domain.Tweet, error)
	FetchAllTweets() ([]domain.Tweet, error)
	FetchAuthorTweets(author_id int) ([]domain.Tweet, error)
	CreateToken(token string) bool
	ReadToken(token string) bool
	UpdateToken(token string, revokeStatus bool) bool
}
