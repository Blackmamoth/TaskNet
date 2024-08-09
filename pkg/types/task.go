package types

import (
	task_model "github.com/blackmamoth/tasknet/pkg/models/task"
	"github.com/blackmamoth/tasknet/pkg/validations"
)

type CreateTaskSchema struct {
	Name          string
	ExecutionMode string
	Priority      string
	UserId        string
}
type TaskService interface {
	GetTaskById(id string) (*task_model.Task, error)
	GetTaskByName(name string) (*task_model.Task, error)
	GetTasks(payload validations.GetTasksSchema, userId string) ([]*task_model.Task, error)
	CreateTask(createTaskSchema CreateTaskSchema) error
	RegisterScriptToTask(scriptId, task_id string) error
}

type TaskRepository interface {
	GetTaskById(id string) (*task_model.Task, error)
	GetTaskByName(name string) (*task_model.Task, error)
	GetTasks(payload validations.GetTasksSchema, userId string) ([]*task_model.Task, error)
	CreateTask(createTaskSchema CreateTaskSchema) error
	RegisterScriptToTask(scriptId, task_id string) error
}
