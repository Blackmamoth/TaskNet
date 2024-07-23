package task_service

import (
	task_model "github.com/blackmamoth/tasknet/pkg/models/task"
	"github.com/blackmamoth/tasknet/pkg/types"
	"github.com/blackmamoth/tasknet/pkg/validations"
)

type Service struct {
	repository types.TaskRepository
}

func New(repository types.TaskRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetTaskById(id string) (*task_model.Task, error) {
	return s.repository.GetTaskById(id)
}

func (s *Service) GetTaskByName(name string) (*task_model.Task, error) {
	return s.repository.GetTaskByName(name)
}

func (s *Service) CreateTask(payload validations.CreateTaskSchema) error {
	return s.repository.CreateTask(payload)
}
