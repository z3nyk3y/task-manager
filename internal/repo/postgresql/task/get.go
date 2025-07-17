package task

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/z3nyk3y/task-manager/internal/models"
)

// FetchTasks - gets provided amount of tasks and set for them status models.Processing.
func (r *Repo) FetchTasks(ctx context.Context, numberOfTasks int) ([]models.Task, error) {
	rows, err := r.db.Query(ctx,
		`WITH chosen_tasks AS (
		    SELECT
				id,
				status_id,
				updated_at,
				created_at
			FROM tasks
			WHERE status_id = $1
			ORDER BY updated_at DESC
			LIMIT $2
		)
		UPDATE tasks
		SET status_id = $3
		WHERE id IN (SELECT id FROM chosen_tasks)
		RETURNING
			id,
			status_id,
			updated_at,
			created_at`, models.New, numberOfTasks, models.Processing)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Task])
	if err != nil {
		return nil, err
	}

	return result, nil
}
