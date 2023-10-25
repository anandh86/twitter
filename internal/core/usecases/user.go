package usecases

import (
	"errors"
	"regexp"
	"strings"

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

func (u userUseCase) PostTweet(body string) (domain.Tweet, error) {
	// business logic here

	// check for validity
	if len(body) > 140 {
		return domain.Tweet{}, errors.New("tweet too long")

	}

	// Replacement illegal words
	{
		// Input string
		input := body

		// Words to be replaced
		wordsToReplace := []string{"kerfuffle", "sharbert", "fornax"}

		// Case-insensitive regular expression pattern
		// The \b word boundary ensures that only full words are matched
		pattern := "(?i)\\b(" + strings.Join(wordsToReplace, "|") + ")\\b"

		// Create the regular expression
		re := regexp.MustCompile(pattern)

		// Replace the words with "****"
		body = re.ReplaceAllString(input, "****")
	}

	tweet := domain.Tweet{Body: body}
	return u.repoImpl.SaveTweet(tweet)
}

func (u userUseCase) GetTweetById(id int) (domain.Tweet, error) {
	return u.repoImpl.GetTweetById(id)
}
