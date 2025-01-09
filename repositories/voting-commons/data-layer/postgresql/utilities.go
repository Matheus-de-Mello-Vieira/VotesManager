package postgresqldatamapper

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresqlConnector struct {
	connectionString string
}

func NewPostgresqlConnector(connectionString string) PostgresqlConnector {
	return PostgresqlConnector{
		connectionString: connectionString,
	}
}

func (connector *PostgresqlConnector) openConnection(ctx context.Context) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(ctx, connector.connectionString)

	if err != nil {
		err = fmt.Errorf("failed to connect to postgresql: %w", err)
	}

	return dbpool, err
}
