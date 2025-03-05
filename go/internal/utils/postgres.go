package utils

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	driver = "pgx"
)

type ParamsPostgresDB struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	AppName  string
}

func NewPostgresDB(params ParamsPostgresDB) (*sqlx.DB, func(), error) {

	dbURL, err := composeURL(params)
	if err != nil {
		return nil, nil, fmt.Errorf("NewPostgresDB: %w", err)
	}

	db, err := sqlx.Open(driver, dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("NewPostgresDB: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, nil, fmt.Errorf("NewPostgresDB: %w", err)
	}
	return db, func() { db.Close() }, nil
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

	if params.AppName == "" {
		params.AppName = "go_service"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&application_name=%s", params.User, params.Password, params.Host, params.Port, params.DBName, params.AppName), nil
}
