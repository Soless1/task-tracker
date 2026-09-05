package json

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"task-tracker/internal/domain/task"
)

func newTestRepository(t *testing.T) *JSONTaskRepository {
	t.Helper()

	path := t.TempDir() + "/tasks.json"

	repo, err := NewJSONTaskRepositoryWithPath(path)
	require.NoError(t, err)

	return repo
}

func TestJSONRepository_Create(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	id := repo.NextID(ctx)

	gotID, err := repo.Create(ctx, task.Task{
		ID:          id,
		Description: "buy milk",
		Status:      task.StatusTodo,
	})

	require.NoError(t, err)
	assert.Equal(t, id, gotID)

	got, err := repo.GetByID(ctx, id)

	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, "buy milk", got.Description)
	assert.Equal(t, task.StatusTodo, got.Status)
}

func TestJSONRepository_GetByID(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	id := repo.NextID(ctx)

	err := func() error {
		_, err := repo.Create(ctx, task.Task{
			ID:          id,
			Description: "write code",
			Status:      task.StatusInProgress,
		})
		return err
	}()

	require.NoError(t, err)

	got, err := repo.GetByID(ctx, id)

	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, "write code", got.Description)
	assert.Equal(t, task.StatusInProgress, got.Status)
}

func TestJSONRepository_GetByID_NotFound(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, 42)

	assert.ErrorIs(t, err, ErrTaskNotFound)
}

func TestJSONRepository_GetAll(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	tasks := []task.Task{
		{
			ID:          repo.NextID(ctx),
			Description: "buy milk",
			Status:      task.StatusTodo,
		},
		{
			ID:          repo.NextID(ctx),
			Description: "write code",
			Status:      task.StatusDone,
		},
	}

	for _, task := range tasks {
		_, err := repo.Create(ctx, task)
		require.NoError(t, err)
	}

	got, err := repo.GetAll(ctx)

	require.NoError(t, err)
	require.Len(t, got, 2)

	assert.ElementsMatch(t, tasks, got)
}

func TestJSONRepository_Update(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	id := repo.NextID(ctx)

	_, err := repo.Create(ctx, task.Task{
		ID:          id,
		Description: "buy milk",
		Status:      task.StatusTodo,
	})
	require.NoError(t, err)

	err = repo.Update(ctx, task.Task{
		ID:          id,
		Description: "buy bread",
		Status:      task.StatusDone,
	})

	require.NoError(t, err)

	got, err := repo.GetByID(ctx, id)

	require.NoError(t, err)

	assert.Equal(t, id, got.ID)
	assert.Equal(t, "buy bread", got.Description)
	assert.Equal(t, task.StatusDone, got.Status)
}

func TestJSONRepository_Update_NotFound(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	err := repo.Update(ctx, task.Task{
		ID:          42,
		Description: "buy bread",
		Status:      task.StatusDone,
	})

	assert.ErrorIs(t, err, ErrTaskNotFound)
}

func TestJSONRepository_Delete(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	id := repo.NextID(ctx)

	_, err := repo.Create(ctx, task.Task{
		ID:          id,
		Description: "buy milk",
		Status:      task.StatusTodo,
	})
	require.NoError(t, err)

	err = repo.Delete(ctx, id)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, id)

	assert.ErrorIs(t, err, ErrTaskNotFound)
}

func TestJSONRepository_Delete_NotFound(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	err := repo.Delete(ctx, 42)

	assert.ErrorIs(t, err, ErrTaskNotFound)
}

func TestJSONRepository_NextID(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	assert.Equal(t, int64(1), repo.NextID(ctx))
	assert.Equal(t, int64(2), repo.NextID(ctx))
	assert.Equal(t, int64(3), repo.NextID(ctx))
}

func TestJSONRepository_NextID_ContinuesFromExistingTasks(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	_, err := repo.Create(ctx, task.Task{
		ID:          5,
		Description: "existing task",
		Status:      task.StatusTodo,
	})
	require.NoError(t, err)

	// Создаём новый repository, чтобы он заново прочитал файл
	repo2, err := NewJSONTaskRepositoryWithPath(repo.path)
	require.NoError(t, err)

	assert.Equal(t, int64(6), repo2.NextID(ctx))
}
