package authentication

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidUserName = errors.New("login failed")
	ErrInvalidPassword = errors.New("email or password is wrong")
)

type UserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserRes struct {
	ID string `json:"ID"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUser(username, password string) *User {
	return &User{
		ID:       uuid.NewString(),
		Username: username,
		Password: password,
	}
}
