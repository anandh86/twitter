package domain

type User struct {
	Email          string
	HashedPassword []byte
	ID             int
}
