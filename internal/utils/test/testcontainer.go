package testh

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	inith "github.com/iv-sukhanov/finance_tracker/internal/utils/init"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type postgresContainer struct {
	testcontainers.Container
}

type postgresContainerOption func(req *testcontainers.ContainerRequest)

func WithWaitStrategy(strats ...wait.Strategy) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.WaitingFor = wait.ForAll(strats...).WithDeadline(time.Minute * 1)
	}
}

func WithPort(port string) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.ExposedPorts = append(req.ExposedPorts, port)
	}
}

func WithInitialDatabase(user, password, dbName string) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.Env["POSTGRES_USER"] = user
		req.Env["POSTGRES_PASSWORD"] = password
		req.Env["POSTGRES_DB"] = dbName
	}
}

func WithHostConfigModigier(modifier func(*container.HostConfig)) func(req *testcontainers.ContainerRequest) {
	return func(req *testcontainers.ContainerRequest) {
		req.HostConfigModifier = modifier
	}
}

func StartConteiner(ctx context.Context, opts ...postgresContainerOption) (*postgresContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres",
		Env:          map[string]string{},
		ExposedPorts: []string{},
		Cmd:          []string{"postgres", "-c", "fsync=off"},
	}

	for _, opt := range opts {
		opt(&req)
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &postgresContainer{container}, nil
}

func NewPGContainer(filenames ...string) (db *sqlx.DB, shut func(), err error) {
	ctx := context.Background()

	const (
		dbName   = "postgres"
		user     = "postgres"
		password = "postgres"
		port     = "5432"
	)

	container, err := StartConteiner(ctx,
		WithPort(port),
		WithInitialDatabase(user, password, dbName),
		WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(time.Second*10)),
		WithHostConfigModigier(func(hc *container.HostConfig) {
			hc.Mounts = filenamesToMounts(filenames...)
		}),
	)
	if err != nil {
		return nil, nil, err
	}

	shut = func() {
		_ = container.Terminate(ctx)
	}

	newPort, err := container.MappedPort(ctx, port)
	if err != nil {
		shut()
		return nil, nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		shut()
		return nil, nil, err
	}

	db, clostDB, err := inith.NewPostgresDB(inith.ParamsPostgresDB{
		User:     user,
		Password: password,
		Host:     host,
		Port:     newPort.Port(),
		DBName:   dbName,
	})
	if err != nil {
		shut()
		return nil, nil, err
	}

	shut = func() {
		clostDB()
		_ = container.Terminate(ctx)
	}

	return db, shut, nil
}

func filenamesToMounts(filenames ...string) []mount.Mount {
	mounts := make([]mount.Mount, 0, len(filenames))

	for i, source := range filenames {
		target := fmt.Sprintf("/docker-entrypoint-initdb.d/%05d.sql", i)

		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: source,
			Target: target,
		})
	}

	return mounts
}
