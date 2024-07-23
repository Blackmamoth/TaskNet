package user_service

import (
	user_model "github.com/blackmamoth/tasknet/pkg/models/user"
	"github.com/blackmamoth/tasknet/pkg/types"
)

type Service struct {
	repository types.UserRepository
}

func New(repository types.UserRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetUserByUsername(username string) (*user_model.User, error) {
	return s.repository.GetUserByUsername(username)
}

func (s *Service) GetUserByEmail(email string) (*user_model.User, error) {
	return s.repository.GetUserByEmail(email)
}

func (s *Service) GetUserById(id string) (*user_model.User, error) {
	return s.repository.GetUserById(id)
}

func (s *Service) CreateUser(user *user_model.User) error {
	return s.repository.CreateUser(user)
}
