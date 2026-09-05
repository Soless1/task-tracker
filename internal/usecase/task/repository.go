package task

import (
	"context"
	"task-tracker/internal/domain/task"
)

type TaskRepository interface {
	Create(ctx context.Context, t task.Task) (int64, error)
	GetByID(ctx context.Context, ID int64) (task.Task, error)
	GetAll(ctx context.Context) ([]task.Task, error)
	Update(ctx context.Context, task task.Task) error
	Delete(ctx context.Context, ID int64) error
	NextID(ctx context.Context) int64
}
