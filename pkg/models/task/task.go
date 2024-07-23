package task_model

import "time"

type TaskStatus string

const (
	Pending    TaskStatus = "pending"
	Scheduled  TaskStatus = "scheduled"
	Processing TaskStatus = "processing"
	Failed     TaskStatus = "failed"
	Processed  TaskStatus = "processed"
)

type Task struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Status     TaskStatus `json:"status"`
	Script     string
	Parameters []string  `json:"parameters"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
