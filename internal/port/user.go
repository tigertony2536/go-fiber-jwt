package port

import "github.com/tigertony2536/go-login/internal/domain"

type UserRepository interface {
	CreateUser(user *domain.UserLogin) error
	GetUsers() ([]*domain.UserLogin, error)
	GetUserByEmail(email string) (*domain.UserLogin, error)
	UpdateUser(user *domain.UserLogin) error
	DeleteUser(id int) error
}

type UserService interface {
}
