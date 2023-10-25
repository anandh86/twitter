package adapters

import (
	"errors"

	"github.com/anandhmaps/chirpy/internal/core/domain"
	"github.com/anandhmaps/chirpy/internal/core/ports"
)

// In memory implementation
func ProvideInMemoryRepo() ports.IRepository {
	return &myInMemoryRepository{
		userMap:           make(map[int]domain.User),
		tweetMap:          make(map[int]domain.Tweet),
		currentNoOfUsers:  0,
		currentNoOfTweets: 0,
	}
}

// myInMemoryRepository implements ports.UserRepository
type myInMemoryRepository struct {
	userMap          map[int]domain.User
	currentNoOfUsers int

	tweetMap          map[int]domain.Tweet
	currentNoOfTweets int
}

func (u *myInMemoryRepository) Save(user domain.User) (domain.User, error) {

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

func (u *myInMemoryRepository) GetUserById(id int) (domain.User, error) {
	user, ok := u.userMap[id]

	if !ok {
		return user, errors.New("user id not found")
	}

	return user, nil
}

func (u *myInMemoryRepository) SaveTweet(tweet domain.Tweet) (domain.Tweet, error) {
	tweetID := u.currentNoOfTweets + 1
	u.currentNoOfTweets = tweetID
	tweet.ID = tweetID
	u.tweetMap[tweetID] = tweet

	return tweet, nil
}

func (u *myInMemoryRepository) GetTweetById(id int) (domain.Tweet, error) {
	tweet, ok := u.tweetMap[id]

	if !ok {
		return tweet, errors.New("tweet id not found")
	}

	return tweet, nil
}
