package validations

type RegisterUserSchema struct {
	Username string `validate:"required,alphanum,min=2,max=20" alias:"username" json:"username"`
	Email    string `validate:"required,email" alias:"email" json:"email"`
	Password string `validate:"required,min=8,max=16" alias:"password" json:"password"`
}

type LoginUserSchema struct {
	Username string `validate:"required,alphanum,min=2,max=20" alias:"username" json:"username"`
	Password string `validate:"required,min=8,max=16" alias:"password" json:"password"`
}

type CreateTaskSchema struct {
	Name          string `validate:"required" alias:"name" json:"name"`
	ExecutionMode string `validate:"required,lowercase,oneof=non-recurring recurring" alias:"execution_mode" json:"execution_mode"`
	Priority      string `validate:"required,lowercase,oneof=low medium high" alias:"priority" json:"priority"`
}

type UploadExecutionScriptSchema struct {
	TaskId     string `validate:"required,uuid" alias:"task_id" json:"task_id"`
	Parameters string `alias:"parameters" json:"parameters"`
}

type GetTasksSchema struct {
	TaskId        string `validate:"omitempty,uuid" alias:"task_id" json:"task_id"`
	ExecutionMode string `validate:"omitempty,lowercase,oneof=non-recurring recurring" alias:"execution_mode" json:"execution_mode"`
	Priority      string `validate:"omitempty,lowercase,oneof=low medium high" alias:"priority" json:"priority"`
	Status        string `validate:"omitempty,lowercase,oneof=unscheduled scheduled active inactive" alias:"status" json:"status"`
	Limit         int    `validate:"omitempty,number" alias:"limit" json:"limit"`
	Offset        string `validate:"omitempty" alias:"offset" json:"offset"`
}
