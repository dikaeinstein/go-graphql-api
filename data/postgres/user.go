package postgres

import (
	"context"
	"fmt"

	"github.com/dikaeinstein/go-graphql-api/data"
	"github.com/pkg/errors"
)

// GetUsersByName retrieves users with name matching the given name
func (p *Postgres) GetUsersByName(ctx context.Context, name string) ([]data.User, error) {
	rows, err := p.QueryContext(ctx, "SELECT * FROM users WHERE name LIKE $1;",
		name)
	if err != nil {
		return nil, errors.Wrap(err, "GetUsersByName failed")
	}
	defer rows.Close()

	users := make([]data.User, 0)
	u := data.User{}
	for rows.Next() {
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age,
			&u.Profession, &u.Friendly)
		if err != nil {
			return users, errors.Wrap(err, "GetUsersByName failed")
		}
		users = append(users, u)
	}

	return users, errors.Wrap(rows.Err(), "GetUsersByName failed")
}

// GetUserByEmail retrieves a single user by email
func (p *Postgres) GetUserByEmail(ctx context.Context, email string) (*data.User, error) {
	row := p.QueryRowContext(ctx, "SELECT * FROM users WHERE email = $1;",
		email)

	var u data.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Age,
		&u.Profession, &u.Friendly)
	if err != nil {
		return nil, errors.Wrap(err, "GetUserByEmail failed")
	}

	return &u, nil
}

// CreateUser creates a new user and returns the ID
func (p *Postgres) CreateUser(ctx context.Context, u data.User) (*data.User, error) {
	query := `
	INSERT INTO users(name, email, age, profession, friendly)
	VALUES($1, $2, $3, $4, $5)
	RETURNING *;`
	row := p.QueryRowContext(ctx, query,
		u.Name, u.Email, u.Age, u.Profession, u.Friendly,
	)

	var newUser data.User
	err := row.Scan(&newUser.ID, &newUser.Name, &newUser.Email, &newUser.Age,
		&newUser.Profession, &newUser.Friendly)
	if err != nil {
		return nil, errors.Wrap(err, "CreateUser failed")
	}

	return &newUser, nil
}

// UpdateUser updates the user that matches `id` with given `payload`
func (p *Postgres) UpdateUser(ctx context.Context, id int,
	payload map[string]interface{}) (*data.User, error) {
	query := prepareUpdateQuery(payload)
	row := p.QueryRowContext(ctx, query, id)

	var u data.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Age,
		&u.Profession, &u.Friendly)
	if err != nil {
		return nil, errors.Wrap(err, "UpdateUser failed")
	}

	return &u, nil
}

// DeleteUser deletes the user that matches `id` from data store
func (p *Postgres) DeleteUser(ctx context.Context, id int) (*data.User, error) {
	row := p.QueryRowContext(ctx,
		"DELETE FROM users WHERE id = $1 RETURNING *;", id)

	var u data.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Age,
		&u.Profession, &u.Friendly)
	if err != nil {
		return nil, errors.Wrap(err, "DeleteUser failed")
	}

	return &u, nil
}

func prepareUpdateQuery(payload map[string]interface{}) string {
	query := "UPDATE users SET"
	for k, v := range payload {
		switch x := v.(type) {
		case int:
			query += fmt.Sprintf(" %s = %d,", k, x)
		case bool:
			query += fmt.Sprintf(" %s = %t,", k, x)
		case string:
			query += fmt.Sprintf(" %s = '%s',", k, x)
		}
	}
	// strip out trailing comma
	query = query[:len(query)-2] + " WHERE id = $1 RETURNING *;"
	return query
}
