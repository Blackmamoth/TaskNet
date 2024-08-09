package task_service

import (
	"fmt"

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

func (s *Service) CreateTask(createTaskSchema types.CreateTaskSchema) error {
	_, err := s.repository.GetTaskByName(createTaskSchema.Name)

	if err == nil {
		return fmt.Errorf("task with name [%s] already exists", createTaskSchema.Name)
	}

	return s.repository.CreateTask(createTaskSchema)
}

func (s *Service) RegisterScriptToTask(scriptId, task_id string) error {
	return s.repository.RegisterScriptToTask(scriptId, task_id)
}

func (s *Service) GetTasks(payload validations.GetTasksSchema, userId string) ([]*task_model.Task, error) {
	return s.repository.GetTasks(payload, userId)
}
