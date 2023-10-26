package usecases

import (
	"errors"
	"regexp"
	"strings"

	"github.com/anandhmaps/chirpy/internal/core/domain"
	"github.com/anandhmaps/chirpy/internal/core/ports"

	"golang.org/x/crypto/bcrypt"
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

func (u userUseCase) CreateUser(emailid string, password string) (domain.User, error) {

	// Generate a salted hash for the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return domain.User{}, errors.ErrUnsupported
	}

	user := domain.User{
		Email:          emailid,
		HashedPassword: hashedPassword,
	}

	savedUser, err1 := u.repoImpl.Save(user)

	if err1 != nil {
		return domain.User{}, errors.ErrUnsupported
	}

	return savedUser, nil
}

func (u userUseCase) UpdateUser(id int, emailid string, password string) error {
	// Generate a salted hash for the password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := domain.User{
		Email:          emailid,
		HashedPassword: hashedPassword,
	}
	return u.repoImpl.UpdateUser(id, user)
}

func (u userUseCase) LoginUser(emailid string, password string) (int, error) {
	// TODO
	userId, err := u.repoImpl.GetUserId(emailid)

	if err != nil {
		return userId, errors.ErrUnsupported
	}

	user, _ := u.repoImpl.GetUserById(userId)

	if !checkPassword(password, user.HashedPassword) {
		return userId, errors.ErrUnsupported
	}

	return userId, err
}

func checkPassword(providedPassword string, hashedPassword []byte) bool {
	// Compare the provided password with the stored hashed password
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(providedPassword))
	return err == nil
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

func (u userUseCase) StoreRefreshToken(token string) bool {
	return u.repoImpl.CreateToken(token)
}

func (u userUseCase) RevokeRefreshToken(token string) bool {

	return u.repoImpl.UpdateToken(token, true)
}

func (u userUseCase) IsRefreshTokenRevoked(token string) bool {
	return u.repoImpl.ReadToken(token)
}
