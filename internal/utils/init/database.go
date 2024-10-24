package inith

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
)

func ConnectDB(ctx context.Context, url string) (db *pgx.Conn, close func(), err error) {

	db, err = pgx.Connect(ctx, url)
	if err != nil {
		return nil, nil, fmt.Errorf("ConnectDB: %w", err)
	}
	return db, func() { db.Close(ctx) }, nil
}

func ConnectDBFromEnv(ctx context.Context) (db *pgx.Conn, close func(), err error) {
	var user, password, host, port, dbName string

	user = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	host = os.Getenv("POSTGRES_HOST")
	port = os.Getenv("POSTGRES_PORT")
	dbName = os.Getenv("POSTGRES_DB")

	url, err := composeURL(user, password, host, port, dbName)
	if err != nil {
		return nil, nil, fmt.Errorf("ConnectDB: %w", err)
	}

	return ConnectDB(ctx, url)

}

func composeURL(user, password, host, port, dbName string) (url string, err error) {

	if user == "" || password == "" || host == "" || port == "" || dbName == "" {
		var missingArgs strings.Builder
		missingArgs.Grow(30) //"user " -> 5, "password " -> 9, "host " -> 5, "port " -> 5, "dbName" -> 6
		if user == "" {
			missingArgs.WriteString("user ")
		}
		if password == "" {
			missingArgs.WriteString("password ")
		}
		if host == "" {
			missingArgs.WriteString("host ")
		}
		if port == "" {
			missingArgs.WriteString("port ")
		}
		if dbName == "" {
			missingArgs.WriteString("dbName")
		}
		return "", fmt.Errorf("composeURL: missing arguments %s", missingArgs.String())
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName), nil
}
