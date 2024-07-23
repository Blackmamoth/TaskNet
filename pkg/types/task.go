package types

import (
	task_model "github.com/blackmamoth/tasknet/pkg/models/task"
	"github.com/blackmamoth/tasknet/pkg/validations"
)

type TaskService interface {
	GetTaskById(id string) (*task_model.Task, error)
	GetTaskByName(name string) (*task_model.Task, error)
	CreateTask(payload validations.CreateTaskSchema) error
}

type TaskRepository interface {
	GetTaskById(id string) (*task_model.Task, error)
	GetTaskByName(name string) (*task_model.Task, error)
	CreateTask(payload validations.CreateTaskSchema) error
}
