package task

import (
	"context"
	"task-tracker/internal/domain/task"
)

func Update(ctx context.Context, repo TaskRepository, ID int64, description string) error {
	var t task.Task
	t, err := repo.GetByID(ctx, ID)
	if err != nil {
		return err
	}
	t.ChangeDescription(description)
	return repo.Update(ctx, t)
}
