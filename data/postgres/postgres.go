package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dikaeinstein/go-graphql-api/data"
)

// Postgres represents the postgres db
type Postgres struct {
	*sql.DB
}

// New opens and returns a postgres DB
func New(connStr string) (*Postgres, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Postgres{db}, nil
}

// GetUsersByName retrieves users with name matching the given name
func (p *Postgres) GetUsersByName(ctx context.Context, name string) ([]data.User, error) {
	rows, err := p.QueryContext(ctx, "SELECT * FROM users WHERE name LIKE $1",
		name)
	if err != nil {
		err = fmt.Errorf("GetUsersByName query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	users := make([]data.User, 0)
	u := data.User{}
	for rows.Next() {
		err := rows.Scan(&u.ID, &u.Name, &u.Age, &u.Profession, &u.Friendly)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}

type connParams struct {
	// Maximum wait for connection, in seconds.
	// Zero or not specified means wait indefinitely.
	connectTimeout int
	dbName         string
	host           string
	password       string
	port           int
	sslMode        string
	user           string
}

// Option sets a connection string param
type Option func(*connParams)

// ConnString returns a connection string based on the parameters it's given
func ConnString(dbName, user string, options ...Option) string {
	// default connection string params
	cp := &connParams{
		connectTimeout: 0,
		dbName:         dbName,
		host:           "localhost",
		port:           5432,
		sslMode:        "require",
		user:           user,
	}

	// apply options to connParams to configure as required
	for _, option := range options {
		option(cp)
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=%s connect_timeout=%d",
		cp.host, cp.port, cp.user, cp.dbName, cp.sslMode, cp.connectTimeout,
	)
}

// Host option sets the host connection string param
func Host(host string) func(*connParams) {
	return func(cp *connParams) {
		cp.host = host
	}
}

// Port option sets the port connection string param
func Port(port int) func(*connParams) {
	return func(cp *connParams) {
		cp.port = port
	}
}

// Password option sets the password connection string param
func Password(password string) func(*connParams) {
	return func(cp *connParams) {
		cp.password = password
	}
}

// SSLMode option sets the sslmode connection string param
func SSLMode(sslMode string) func(*connParams) {
	return func(cp *connParams) {
		cp.sslMode = sslMode
	}
}

// ConnectTimeout option sets the connect_timeout connection string param
func ConnectTimeout(timeout int) func(*connParams) {
	return func(cp *connParams) {
		cp.connectTimeout = timeout
	}
}
