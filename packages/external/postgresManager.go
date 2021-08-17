package external

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type IPostgresManager interface {
	Connect(connectionString string) (err error)
	DeleteDatabase(databse string) (err error)
	Close()
}

type PostgresManager struct {
	connection *pgx.Conn
	ctx        context.Context
}

func NewPostgresManager() *PostgresManager {
	return &PostgresManager{
		ctx: context.Background(),
	}
}

func (p *PostgresManager) Close() {
	defer p.connection.Close(p.ctx)
}

func (p *PostgresManager) Connect(connectionString string) (err error) {
	conn, err := pgx.Connect(p.ctx, connectionString)

	if err != nil {
		return errors.Wrap(err, "Failed to connect to database")
	}
	p.connection = conn
	return
}

func (p *PostgresManager) DeleteDatabase(database string) (err error) {
	closeDbConnections := `SELECT pg_terminate_backend(pg_stat_activity.pid)
	FROM pg_stat_activity
	WHERE pg_stat_activity.datname = '%s'
		and pid <> pg_backend_pid();`
	deleteDb := `DROP DATABASE IF EXISTS "%s"`

	_, err = p.connection.Exec(p.ctx, fmt.Sprintf(closeDbConnections, database))

	if err != nil {
		return errors.Wrapf(err, "Failed to close database connections for %s", database)
	}

	_, err = p.connection.Exec(p.ctx, fmt.Sprintf(deleteDb, database))

	if err != nil {
		return errors.Wrapf(err, "Failed to delete database %s", database)
	}

	return
}
