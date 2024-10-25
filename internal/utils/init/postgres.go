package inith

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ParamsPostgresDB struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

func NewPostgresDB(params ParamsPostgresDB) (*sqlx.DB, error) {

	dbURL, err := composeURL(params)
	if err != nil {
		return nil, fmt.Errorf("NewPostgresDB: %w", err)
	}

	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("NewPostgresDB: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("NewPostgresDB: %w", err)
	}
	return db, nil
}

func composeURL(params ParamsPostgresDB) (url string, err error) {

	if params.User == "" || params.Password == "" || params.Host == "" || params.Port == "" || params.DBName == "" {
		var missingArgs strings.Builder
		missingArgs.Grow(30) //"user " -> 5, "password " -> 9, "host " -> 5, "port " -> 5, "dbName" -> 6
		if params.User == "" {
			missingArgs.WriteString("user ")
		}
		if params.Password == "" {
			missingArgs.WriteString("password ")
		}
		if params.Host == "" {
			missingArgs.WriteString("host ")
		}
		if params.Port == "" {
			missingArgs.WriteString("port ")
		}
		if params.DBName == "" {
			missingArgs.WriteString("dbName")
		}
		return "", fmt.Errorf("composeURL: missing arguments %s", missingArgs.String())
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", params.User, params.Password, params.Host, params.Port, params.DBName), nil
}
