package task

import (
	"context"
	"iter"
	"slices"
	"task-tracker/internal/domain/task"
)

func List(ctx context.Context, repo TaskRepository, status *task.Status) (iter.Seq[task.Task], error) {
	tasks, err := repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if status == nil {
		return slices.Values(tasks), nil
	}

	return func(yield func(task.Task) bool) {
		for _, t := range tasks {
			if t.Status == *status {
				if !yield(t) {
					return
				}
			}
		}
	}, nil
}
