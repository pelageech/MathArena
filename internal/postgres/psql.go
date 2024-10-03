package postgres

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PSQLDatabase accesses to PostgresSQL with other fields.
type PSQLDatabase struct {
	*pgxpool.Pool
	logger *log.Logger
}

// NewPSQLDatabase creates a new connection to db and tries to connect to it.
func NewPSQLDatabase(ctx context.Context, connString string, logger *log.Logger) (*PSQLDatabase, error) {
	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	return &PSQLDatabase{
		Pool:   conn,
		logger: logger,
	}, nil
}

func (d *PSQLDatabase) Logger() *log.Logger {
	return d.logger
}
