package authentication

import (
	"log"
	"sync"
)

type Authentication interface {
	Login(username string, password string) error
	CreateUser(userReq *UserReq) (User, error)
}

type AuthService struct {
	Users map[string]*User
	mu    sync.Mutex
}

var _ Authentication = (*AuthService)(nil)

func NewAuthService() Authentication {
	userMap := make(map[string]*User)
	testUser := NewUser("preeti", "admin123")
	userMap[testUser.ID] = testUser

	return &AuthService{
		Users: userMap,
	}
}

func (ul *AuthService) Login(username string, password string) error {
	if username == "" || password == "" {
		return ErrInvalidPassword
	}

	ul.mu.Lock()
	defer ul.mu.Unlock()

	for _, user := range ul.Users {
		if user.Username == username && user.Password == password {
			log.Printf("user login successful: %s", username)
			return nil
		}
	}

	return ErrInvalidPassword
}

func (ul *AuthService) CreateUser(userReq *UserReq) (User, error) {
	if len(userReq.Password) <= 3 {
		log.Printf("Minimum length than threshold %v", len(userReq.Password))
		return User{}, ErrInvalidPassword
	}

	ul.mu.Lock()
	defer ul.mu.Unlock()

	newUser := NewUser(userReq.Username, userReq.Password)
	ul.Users[newUser.ID] = newUser
	return *newUser, nil
}
