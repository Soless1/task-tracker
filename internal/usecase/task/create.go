package task

import (
	"context"
	"task-tracker/internal/domain/task"
	"time"
)

func CreateTask(ctx context.Context, repo TaskRepository, description string) (int64, error) {
	t := task.CreateTask(description, time.Now(), repo.NextID(ctx))
	return repo.Create(ctx, t)

}
