package task_repository

import (
	"context"
	"fmt"

	task_model "github.com/blackmamoth/tasknet/pkg/models/task"
	"github.com/blackmamoth/tasknet/pkg/types"
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
	args := pgx.NamedArgs{
		"id": id,
	}
	rows, err := r.conn.Query(context.Background(), "SELECT * FROM task WHERE id = @id", args)

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
		return nil, fmt.Errorf("task [%s] does not exist", id)
	}

	return t, nil
}

func (r *Repository) GetTaskByName(name string) (*task_model.Task, error) {
	args := pgx.NamedArgs{
		"name": name,
	}
	rows, err := r.conn.Query(context.Background(), "SELECT * FROM task WHERE name = @name", args)

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

func (r *Repository) CreateTask(createTaskSchema types.CreateTaskSchema) error {
	args := pgx.NamedArgs{
		"name":           createTaskSchema.Name,
		"execution_mode": createTaskSchema.ExecutionMode,
		"priority":       createTaskSchema.Priority,
		"user_id":        createTaskSchema.UserId,
	}

	_, err := r.conn.Exec(context.Background(), "INSERT INTO task(name, execution_mode, priority, user_id) VALUES(@name, @execution_mode, @priority, @user_id)", args)
	return err
}

func (r *Repository) RegisterScriptToTask(scriptId, task_id string) error {
	args := pgx.NamedArgs{
		"script_id": scriptId,
		"id":        task_id,
	}
	_, err := r.conn.Exec(context.Background(), "UPDATE task SET script_id = @script_id WHERE id = @id", args)
	return err
}

func (r *Repository) GetTasks(payload validations.GetTasksSchema, userId string) ([]*task_model.Task, error) {
	args := pgx.NamedArgs{
		"user_id": userId,
	}
	sqlString := "SELECT * FROM task WHERE user_id = @user_id"
	if payload.TaskId != "" {
		args["id"] = payload.TaskId
		sqlString += " AND id = @id"
	}

	if payload.ExecutionMode != "" {
		args["execution_mode"] = payload.ExecutionMode
		sqlString += " AND execution_mode = @execution_mode"
	}

	if payload.Priority != "" {
		args["priority"] = payload.Priority
		sqlString += " AND priority = @priority"
	}

	if payload.Status != "" {
		args["status"] = payload.Status
		sqlString += " AND status = @status"
	}

	if payload.Offset != "" {
		args["created_at"] = payload.Offset
		sqlString += " AND created_at > @created_at ORDER BY created_at"
	}

	if payload.Limit != 0 {
		args["limit"] = payload.Limit
		sqlString += " LIMIT @limit"
	}

	rows, err := r.conn.Query(context.Background(), sqlString, args)

	if err != nil {
		return nil, err
	}

	tasks := []*task_model.Task{}

	for rows.Next() {
		t := new(task_model.Task)
		t, err = scanRows(rows)

		if err != nil {
			return nil, err
		}

		if t.Id == "" {
			return nil, fmt.Errorf("an error occured while fetching tasks")
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func scanRows(rows pgx.Rows) (*task_model.Task, error) {
	t := new(task_model.Task)

	err := rows.Scan(
		&t.Id,
		&t.Name,
		&t.Status,
		&t.ExecutionMode,
		&t.Priority,
		&t.PreviousRun,
		&t.NextRun,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.ScriptId,
		&t.UserId,
	)

	return t, err
}
