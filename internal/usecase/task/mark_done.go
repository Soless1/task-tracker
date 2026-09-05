package task

import (
	"context"
	"task-tracker/internal/domain/task"
)

func MarkDone(ctx context.Context, repo TaskRepository, ID int64) error {
	var t task.Task
	t, err := repo.GetByID(ctx, ID)
	if err != nil {
		return err
	}
	t.MarkDone()
	return repo.Update(ctx, t)
}
