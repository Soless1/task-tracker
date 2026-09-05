package task

import (
	"context"
	"task-tracker/internal/domain/task"
)

func MarkInProgress(ctx context.Context, repo TaskRepository, ID int64) error {
	var t task.Task
	t, err := repo.GetByID(ctx, ID)
	if err != nil {
		return err
	}
	t.MarkInProgress()
	return repo.Update(ctx, t)
}
