package port

import "github.com/tigertony2536/go-login/internal/core/domain"

type UserRepository interface {
	CreateUser(user *domain.UserLogin) error
	GetUsers() ([]*domain.UserLogin, error)
	GetUserByID(is int) (*domain.UserLogin, error)
	GetUserByEmail(email string) (*domain.UserLogin, error)
	UpdateUser(user *domain.UserLogin) error
	DeleteUser(id int) error
}

type AuthService interface {
	//user
	Login()
	Refresh()
	Delete()
	Update()
	//admin

}
