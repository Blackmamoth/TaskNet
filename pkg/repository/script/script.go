package script_repository

import (
	"context"
	"fmt"

	script_model "github.com/blackmamoth/tasknet/pkg/models/script"
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

func (r *Repository) CreateScript(script *script_model.Script) error {
	args := pgx.NamedArgs{
		"name":          script.Name,
		"original_name": script.OriginalName,
		"parameters":    script.Parameters,
	}

	_, err := r.conn.Exec(context.Background(), "INSERT INTO script(name, original_name, parameters) VALUES(@name, @original_name, @parameters);", args)
	return err
}

func (r *Repository) GetScriptByName(name string) (*script_model.Script, error) {
	args := pgx.NamedArgs{
		"name": name,
	}

	rows, err := r.conn.Query(context.Background(), "SELECT * FROM script WHERE name = @name", args)

	if err != nil {
		return nil, err
	}

	s := new(script_model.Script)

	for rows.Next() {
		s, err = scanRows(rows)

		if err != nil {
			return nil, err
		}
	}

	if s.Id == "" {
		return nil, fmt.Errorf("script with name [%s] does not exist", name)
	}

	return s, nil
}

func scanRows(rows pgx.Rows) (*script_model.Script, error) {
	t := new(script_model.Script)

	err := rows.Scan(
		&t.Id,
		&t.Name,
		&t.OriginalName,
		&t.Parameters,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	return t, err
}
