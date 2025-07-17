package task

import (
	"context"

	"github.com/z3nyk3y/task-manager/internal/models"
)

// UpdateTask - updates task
func (r *Repo) UpdateTask(ctx context.Context, task models.Task) error {
	_, err := r.db.Exec(ctx,
		`UPDATE tasks SET
			status_id = $1
		WHERE id = $2`, task.Status, task.Id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateTasksStatus updates common status for provided tasks.
func (r *Repo) UpdateTasksStatus(ctx context.Context, tasks []models.Task, status models.TaskStatus) error {
	var ids = make([]int64, 0, len(tasks))

	for _, task := range tasks {
		ids = append(ids, task.Id)
	}

	_, err := r.db.Exec(ctx,
		`UPDATE tasks SET
			status_id = $1
		WHERE id = ANY($2)`, status, ids)
	if err != nil {
		return err
	}

	return nil
}
