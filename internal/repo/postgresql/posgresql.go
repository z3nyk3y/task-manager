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
	pool, err := pgxpool.New(ctx, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?pool_max_conns=30&pool_min_conns=5", cfg.Login, cfg.Password, cfg.Host, cfg.Port, cfg.Name))
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func (r *Repos) Close(pool *pgxpool.Pool) {
	pool.Close()
}
