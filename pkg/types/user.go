package types

import user_model "github.com/blackmamoth/tasknet/pkg/models/user"

type UserService interface {
	GetUserByUsername(username string) (*user_model.User, error)
	GetUserByEmail(email string) (*user_model.User, error)
	GetUserById(id string) (*user_model.User, error)
	CreateUser(user *user_model.User) error
}

type UserRepository interface {
	GetUserByUsername(username string) (*user_model.User, error)
	GetUserByEmail(email string) (*user_model.User, error)
	GetUserById(id string) (*user_model.User, error)
	CreateUser(user *user_model.User) error
}
