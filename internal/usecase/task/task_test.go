package task

import (
	"context"
	"errors"
	"task-tracker/internal/domain/task"
	"task-tracker/internal/usecase/task/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockTaskRepository(ctrl)

	repo.EXPECT().
		Create(
			gomock.Any(),
			gomock.Cond(func(t task.Task) bool {
				return t.Description == "buy milk" &&
					t.Status == task.StatusTodo
			}),
		).
		Return(int64(42), nil)

	got, err := CreateTask(
		context.Background(),
		repo,
		"buy milk",
	)

	assert.NoError(t, err)
	assert.Equal(t, int64(42), got)
}

func TestCreateTask_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockTaskRepository(ctrl)

	expectedErr := errors.New("database error")

	repo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(int64(0), expectedErr)

	got, err := CreateTask(
		context.Background(),
		repo,
		"buy milk",
	)

	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, int64(0), got)
}

func TestCreateTask_TaskData(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockTaskRepository(ctrl)

	repo.EXPECT().
		Create(
			gomock.Any(),
			gomock.Cond(func(t task.Task) bool {
				return t.Description == "buy milk" &&
					t.Status == task.StatusTodo &&
					t.ID == 0 &&
					!t.CreatedAt.IsZero() &&
					!t.UpdatedAt.IsZero()
			}),
		).
		Return(int64(42), nil)

	got, err := CreateTask(
		context.Background(),
		repo,
		"buy milk",
	)

	assert.NoError(t, err)
	assert.Equal(t, int64(42), got)
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockTaskRepository(ctrl)

	repo.EXPECT().
		GetByID(gomock.Any(), int64(42)).
		Return(task.Task{
			ID:          42,
			Description: "old description",
			Status:      task.StatusTodo,
		}, nil)

	repo.EXPECT().
		Update(
			gomock.Any(),
			gomock.Cond(func(t task.Task) bool {
				return t.ID == 42 &&
					t.Description == "buy milk" &&
					t.Status == task.StatusTodo
			}),
		).
		Return(nil)

	err := Update(
		context.Background(),
		repo,
		42,
		"buy milk",
	)

	assert.NoError(t, err)
}

func TestUpdate_GetByIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockTaskRepository(ctrl)

	expectedErr := errors.New("task not found")

	repo.EXPECT().
		GetByID(gomock.Any(), int64(42)).
		Return(task.Task{}, expectedErr)

	err := Update(
		context.Background(),
		repo,
		42,
		"buy milk",
	)

	assert.ErrorIs(t, err, expectedErr)
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockTaskRepository(ctrl)

	repo.EXPECT().
		Delete(gomock.Any(), int64(42)).
		Return(nil)

	err := Delete(
		context.Background(),
		repo,
		42,
	)

	assert.NoError(t, err)
}

func TestDelete_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockTaskRepository(ctrl)

	expectedErr := errors.New("database error")

	repo.EXPECT().
		Delete(gomock.Any(), int64(42)).
		Return(expectedErr)

	err := Delete(
		context.Background(),
		repo,
		42,
	)

	assert.ErrorIs(t, err, expectedErr)
}

func TestMarkDone(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockTaskRepository(ctrl)

	repo.EXPECT().
		GetByID(gomock.Any(), int64(42)).
		Return(task.Task{
			ID:          42,
			Description: "buy milk",
			Status:      task.StatusTodo,
		}, nil)

	repo.EXPECT().
		Update(
			gomock.Any(),
			gomock.Cond(func(t task.Task) bool {
				return t.ID == 42 &&
					t.Status == task.StatusDone
			}),
		).
		Return(nil)

	err := MarkDone(
		context.Background(),
		repo,
		42,
	)

	assert.NoError(t, err)
}

func TestMarkDone_GetByIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockTaskRepository(ctrl)

	expectedErr := errors.New("task not found")

	repo.EXPECT().
		GetByID(gomock.Any(), int64(42)).
		Return(task.Task{}, expectedErr)

	err := MarkDone(
		context.Background(),
		repo,
		42,
	)

	assert.ErrorIs(t, err, expectedErr)
}
func TestList(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockTaskRepository(ctrl)

	expected := []task.Task{
		{
			ID:          1,
			Description: "buy milk",
			Status:      task.StatusTodo,
		},
		{
			ID:          2,
			Description: "write code",
			Status:      task.StatusDone,
		},
	}

	repo.EXPECT().
		GetAll(gomock.Any()).
		Return(expected, nil)

	seq, err := List(
		context.Background(),
		repo,
		nil,
	)

	assert.NoError(t, err)

	var got []task.Task
	for t := range seq {
		got = append(got, t)
	}

	assert.Equal(t, expected, got)
}

func TestList_ByStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockTaskRepository(ctrl)

	status := task.StatusDone

	allTasks := []task.Task{
		{
			ID:          1,
			Description: "buy milk",
			Status:      task.StatusTodo,
		},
		{
			ID:          2,
			Description: "write code",
			Status:      task.StatusDone,
		},
		{
			ID:          3,
			Description: "read book",
			Status:      task.StatusDone,
		},
	}

	repo.EXPECT().
		GetAll(gomock.Any()).
		Return(allTasks, nil)

	seq, err := List(context.Background(), repo, &status)

	assert.NoError(t, err)

	var got []task.Task
	for t := range seq {
		got = append(got, t)
	}

	expected := []task.Task{
		allTasks[1],
		allTasks[2],
	}

	assert.Equal(t, expected, got)
}
