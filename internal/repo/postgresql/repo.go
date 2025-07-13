package postgresql

import (
	"gitlab.com/eightlix-group/task-manager/internal/repo/postgresql/task"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repos struct {
	Task *task.Repo
}

func NewRepository(db *pgxpool.Pool) *Repos {
	return &Repos{
		Task: task.NewRepo(db),
	}
}
