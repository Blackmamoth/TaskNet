package user_repository

import (
	"context"
	"fmt"

	user_model "github.com/blackmamoth/tasknet/pkg/models/user"
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

func (r *Repository) GetUserByUsername(username string) (*user_model.User, error) {
	args := pgx.NamedArgs{
		"username": username,
	}

	rows, err := r.conn.Query(context.Background(), "SELECT * FROM \"user\" WHERE username = @username LIMIT 1", args)

	if err != nil {
		return nil, err
	}

	u := new(user_model.User)

	for rows.Next() {
		u, err = scanRows(rows)

		if err != nil {
			return nil, err
		}
	}

	if u.Id == "" {
		return nil, fmt.Errorf("user with username [%s] does not exist", username)
	}

	return u, nil
}

func (r *Repository) GetUserByEmail(email string) (*user_model.User, error) {
	args := pgx.NamedArgs{
		"email": email,
	}

	rows, err := r.conn.Query(context.Background(), "SELECT * FROM \"user\" WHERE email = @email LIMIT 1", args)

	if err != nil {
		return nil, err
	}

	u := new(user_model.User)

	for rows.Next() {
		u, err = scanRows(rows)

		if err != nil {
			return nil, err
		}
	}

	if u.Id == "" {
		return nil, fmt.Errorf("user with email [%s] does not exist", email)
	}

	return u, nil
}

func (r *Repository) GetUserById(id string) (*user_model.User, error) {
	args := pgx.NamedArgs{
		"id": id,
	}

	rows, err := r.conn.Query(context.Background(), "SELECT * FROM \"user\" WHERE id = @id LIMIT 1", args)

	if err != nil {
		return nil, err
	}

	u := new(user_model.User)

	for rows.Next() {
		u, err = scanRows(rows)

		if err != nil {
			return nil, err
		}
	}

	if u.Id == "" {
		return nil, fmt.Errorf("user [%s] does not exist", id)
	}

	return u, nil
}

func (r *Repository) CreateUser(user *user_model.User) error {
	args := pgx.NamedArgs{
		"username": user.Username,
		"email":    user.Email,
		"password": user.Password,
	}

	_, err := r.conn.Exec(context.Background(), "INSERT INTO \"user\"(username, email, password) VALUES(@username, @email, @password)", args)
	return err
}

func scanRows(rows pgx.Rows) (*user_model.User, error) {
	u := new(user_model.User)

	err := rows.Scan(
		&u.Id,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	return u, err
}
