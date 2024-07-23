package validations

type CreateTaskSchema struct {
	Name       string   `validate:"required,min=4,max=34" json:"name"`
	Parameters []string `validate:"dive" json:"parameters"`
}

type RegisterUserSchema struct {
	Username string `validate:"required,min=2,max=20" json:"username"`
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required,min=8,max=16" json:"password"`
}

type LoginUserSchema struct {
	Username string `validate:"required,min=2,max=20" json:"username"`
	Password string `validate:"required,min=8,max=16" json:"password"`
}
