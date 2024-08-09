package task_model

import (
	"database/sql"
	"time"
)

type TaskStatus string
type TaskType string
type TaskExecutionMode string
type TaskPriority string

const (
	UnScheduled TaskStatus = "unscheduled"
	Scheduled   TaskStatus = "scheduled"
	Active      TaskStatus = "active"
	inactive    TaskStatus = "inactive"
)

const (
	NonRecurring TaskExecutionMode = "non-recurring"
	Recurring    TaskExecutionMode = "recurring"
)

const (
	Low    TaskPriority = "low"
	Medium TaskPriority = "medium"
	High   TaskPriority = "high"
)

type Task struct {
	Id            string            `json:"id"`
	Name          string            `json:"name"`
	ScriptId      sql.NullString    `json:"script_id"`
	UserId        string            `json:"user_id"`
	Status        TaskStatus        `json:"status"`
	ExecutionMode TaskExecutionMode `json:"execution_mode"`
	Priority      TaskPriority      `json:"priority"`
	PreviousRun   *time.Time        `json:"previous_run"`
	NextRun       *time.Time        `json:"next_run"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}
