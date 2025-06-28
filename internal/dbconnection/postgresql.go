package dbconnection

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	timeout = time.Duration(time.Second * 5)
)

func NewPostgreSQLConnectionPool(ctx context.Context, connectionURI string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connectionURI)
	if err != nil {
		return nil, fmt.Errorf("failed to create new db conn pool: %w", err)
	}

	internalCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err = pool.Ping(internalCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return pool, nil
}
