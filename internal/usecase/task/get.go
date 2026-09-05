package task

import (
	"context"
	"task-tracker/internal/domain/task"
)

func GetTask(ctx context.Context, repo TaskRepository, id int64) (task.Task, error) {
	return repo.GetByID(ctx, id)
}
