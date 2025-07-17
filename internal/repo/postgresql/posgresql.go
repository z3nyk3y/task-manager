package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Login    string
	Password string
	Host     string
	Port     string
	Name     string
}

func NewPostgreSqlDB(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?pool_max_conns=30&pool_min_conns=5&sslmode=disable",
		cfg.Login, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func (r *Repos) Close(pool *pgxpool.Pool) {
	pool.Close()
}
