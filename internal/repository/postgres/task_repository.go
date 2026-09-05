package postgres

import (
	"context"
	"task-tracker/internal/domain/task"
	"task-tracker/internal/repository/postgres/generated"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	queries *generated.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		queries: generated.New(pool),
	}
}

func (r *Repository) Create(ctx context.Context, t task.Task) (int64, error) {
	return r.queries.CreateTask(ctx, generated.CreateTaskParams{
		Description: t.Description,
		Status:      toDBStatus(t.Status),
	})
}

func (r *Repository) GetByID(ctx context.Context, ID int64) (task.Task, error) {
	t, err := r.queries.GetTask(ctx, ID)
	if err != nil {
		return task.Task{}, err
	}

	return fromDBTask(t), nil
}

func (r *Repository) GetAll(ctx context.Context) ([]task.Task, error) {
	tasks, err := r.queries.GetTasks(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]task.Task, 0, len(tasks))

	for _, t := range tasks {
		result = append(result, fromDBTask(t))
	}

	return result, nil
}

func (r *Repository) Update(ctx context.Context, t task.Task) error {
	return r.queries.UpdateTask(ctx, generated.UpdateTaskParams{
		ID:          t.ID,
		Description: t.Description,
		Status:      toDBStatus(t.Status),
	})
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	return r.queries.DeleteTask(ctx, id)
}

func (r *Repository) NextID(ctx context.Context) int64 {
	return 0
}

func toDBStatus(status task.Status) generated.TaskTrackerTaskStatus {
	switch status {
	case task.StatusTodo:
		return generated.TaskTrackerTaskStatusTodo
	case task.StatusInProgress:
		return generated.TaskTrackerTaskStatusInProgress
	case task.StatusDone:
		return generated.TaskTrackerTaskStatusDone
	case task.StatusCanceled:
		return generated.TaskTrackerTaskStatusCancelled
	default:
		panic("unknown task status")
	}
}
func fromDBStatus(status generated.TaskTrackerTaskStatus) task.Status {
	switch status {
	case generated.TaskTrackerTaskStatusTodo:
		return task.StatusTodo
	case generated.TaskTrackerTaskStatusInProgress:
		return task.StatusInProgress
	case generated.TaskTrackerTaskStatusDone:
		return task.StatusDone
	case generated.TaskTrackerTaskStatusCancelled:
		return task.StatusCanceled
	default:
		panic("unknown task status")
	}
}

func fromDBTask(t generated.TaskTrackerTask) task.Task {
	return task.Task{
		ID:          t.ID,
		Description: t.Description,
		Status:      fromDBStatus(t.Status),
		CreatedAt:   t.CreatedAt.Time,
		UpdatedAt:   t.UpdatedAt.Time,
	}
}
