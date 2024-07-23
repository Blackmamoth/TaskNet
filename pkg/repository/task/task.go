package task_repository

import (
	"context"
	"fmt"

	task_model "github.com/blackmamoth/tasknet/pkg/models/task"
	"github.com/blackmamoth/tasknet/pkg/validations"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *Repository {
	return &Repository{
		conn: conn,
	}
}

func (r *Repository) GetTaskById(id string) (*task_model.Task, error) {
	return nil, nil
}

func (r *Repository) GetTaskByName(name string) (*task_model.Task, error) {
	args := pgx.NamedArgs{
		"name": name,
	}
	rows, err := r.conn.Query(context.Background(), "SELECT * FROM tasks WHERE name = @name", args)

	if err != nil {
		return nil, err
	}

	t := new(task_model.Task)

	for rows.Next() {
		t, err = scanRows(rows)

		if err != nil {
			return nil, err
		}
	}

	if t.Id == "" {
		return nil, fmt.Errorf("task with name [%s] does not exist", name)
	}

	return t, nil
}

func (r *Repository) createTask(payload validations.CreateTaskSchema) error {
	args := pgx.NamedArgs{
		"name":       payload.Name,
		"parameters": payload.Parameters,
		"script":     "ls",
	}

	_, err := r.conn.Exec(context.Background(), "INSERT INTO tasks(name, parameters, script) VALUES(@name, @parameters, @script)", args)
	return err
}

func (r *Repository) CreateTask(payload validations.CreateTaskSchema) error {
	_, err := r.GetTaskByName(payload.Name)

	if err == nil {
		return fmt.Errorf("task with name [%s] already exists", payload.Name)
	}

	return r.createTask(payload)
}

func scanRows(rows pgx.Rows) (*task_model.Task, error) {
	t := new(task_model.Task)

	err := rows.Scan(
		&t.Id,
		&t.Name,
		&t.Status,
		&t.Script,
		&t.Parameters,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	return t, err
}
