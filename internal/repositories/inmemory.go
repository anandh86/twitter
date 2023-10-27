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
		emaild2idMap:      make(map[string]int),
		tokenRepo:         make(map[string]bool),
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

	emaild2idMap map[string]int

	tokenRepo map[string]bool
}

func (u *myInMemoryRepository) CreateToken(token string) bool {

	if _, ok := u.tokenRepo[token]; ok {
		return false
	}

	u.tokenRepo[token] = false
	return true
}

func (u *myInMemoryRepository) ReadToken(token string) bool {

	return u.tokenRepo[token]
}

func (u *myInMemoryRepository) UpdateToken(token string, revokeStatus bool) bool {

	if _, ok := u.tokenRepo[token]; !ok {
		return false
	}

	u.tokenRepo[token] = revokeStatus

	return true
}

func (u *myInMemoryRepository) Save(user domain.User) (domain.User, error) {

	if _, ok := u.emaild2idMap[user.Email]; ok {
		// user already present
		return user, errors.ErrUnsupported
	}

	userId := u.currentNoOfUsers + 1
	u.currentNoOfUsers = userId
	user.ID = userId

	u.userMap[userId] = user
	u.emaild2idMap[user.Email] = userId

	return user, nil
}

func (u *myInMemoryRepository) UpdateUserMembership(id int, isMember bool) error {
	// for valid item, update the data structures
	dbUser, ok := u.userMap[id]

	if !ok {
		return errors.ErrUnsupported
	}

	dbUser.IsChirpyRed = isMember
	u.userMap[id] = dbUser
	return nil
}

func (u *myInMemoryRepository) UpdateUser(id int, user domain.User) error {
	// for valid item, update the data structures
	dbUser, ok := u.userMap[id]

	if !ok {
		return errors.ErrUnsupported
	}

	if user.Email != "" && user.Email != dbUser.Email {
		// delete old email
		delete(u.emaild2idMap, dbUser.Email)
		dbUser.Email = user.Email
		u.emaild2idMap[user.Email] = id
	}

	user.ID = id
	u.userMap[id] = user
	return nil
}

func (u *myInMemoryRepository) GetUserById(id int) (domain.User, error) {
	user, ok := u.userMap[id]

	if !ok {
		return user, errors.New("user id not found")
	}

	return user, nil
}

func (u *myInMemoryRepository) GetUserId(emailid string) (int, error) {

	userId, ok := u.emaild2idMap[emailid]

	if !ok {
		// user not present
		return 0, errors.ErrUnsupported
	}

	return userId, nil
}

func (u *myInMemoryRepository) SaveTweet(tweet domain.Tweet) (domain.Tweet, error) {
	tweetID := u.currentNoOfTweets + 1
	u.currentNoOfTweets = tweetID
	tweet.ID = tweetID
	u.tweetMap[tweetID] = tweet

	return tweet, nil
}

func (u *myInMemoryRepository) DeleteTweet(tweet domain.Tweet) error {
	tweetID := tweet.ID

	if _, ok := u.tweetMap[tweetID]; !ok {
		// tweet not present
		return errors.ErrUnsupported
	}

	delete(u.tweetMap, tweetID)

	return nil
}

func (u *myInMemoryRepository) GetTweetById(id int) (domain.Tweet, error) {
	tweet, ok := u.tweetMap[id]

	if !ok {
		return tweet, errors.New("tweet id not found")
	}

	return tweet, nil
}

func (u *myInMemoryRepository) FetchAllTweets() ([]domain.Tweet, error) {
	tweets := make([]domain.Tweet, 0)

	for _, tweet := range u.tweetMap {
		tweets = append(tweets, tweet)
	}

	return tweets, nil
}

func (u *myInMemoryRepository) FetchAuthorTweets(author_id int) ([]domain.Tweet, error) {

	tweets := make([]domain.Tweet, 0)

	for _, tweet := range u.tweetMap {
		if tweet.Author != author_id {
			continue
		}
		tweets = append(tweets, tweet)
	}

	return tweets, nil
}
