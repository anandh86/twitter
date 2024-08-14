package ports

import "github.com/anandh86/chirpy/internal/core/domain"

// IUseCase is a primary port that the core must respond to
type IUseCase interface {
	CreateUser(emailid string, password string) (domain.User, error)
	UpdateUser(id int, emailid string, password string) error
	UpdateUserMembership(id int, isMember bool) error
	GetUserById(id int) (domain.User, error)
	LoginUser(emailid string, password string) (int, error)
	PostTweet(body string, author_id int) (domain.Tweet, error)
	DeleteTweet(tweetId int, author_id int) error
	GetTweetById(id int) (domain.Tweet, error)
	GetAllTweets() ([]domain.Tweet, error)
	GetAuthorTweets(author_id int) ([]domain.Tweet, error)
	StoreRefreshToken(token string) bool
	RevokeRefreshToken(token string) bool
	IsRefreshTokenRevoked(token string) bool
}
