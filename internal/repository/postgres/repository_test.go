package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"task-tracker/internal/domain/task"
	"task-tracker/migrations"
)

const testDatabaseURL = "postgres://postgres:postgres@localhost:5432/task_tracker"

var testPool *pgxpool.Pool
var testRepo *Repository

func TestMain(m *testing.M) {
	ctx := context.Background()

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = testDatabaseURL
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	if err := migrations.Run(pool); err != nil {
		panic(err)
	}

	testPool = pool
	testRepo = NewRepository(pool)

	os.Exit(m.Run())
}

func cleanDB(t *testing.T) {
	t.Helper()

	_, err := testPool.Exec(
		context.Background(),
		"TRUNCATE task_tracker.tasks RESTART IDENTITY",
	)
	require.NoError(t, err)
}

func createTestTask(t *testing.T, description string, status task.Status) int64 {
	t.Helper()

	id, err := testRepo.Create(
		context.Background(),
		task.Task{
			Description: description,
			Status:      status,
		},
	)

	require.NoError(t, err)

	return id
}

func TestRepository_Create(t *testing.T) {
	cleanDB(t)

	ctx := context.Background()

	id, err := testRepo.Create(ctx, task.Task{
		Description: "buy milk",
		Status:      task.StatusTodo,
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1), id)
}

func TestRepository_GetByID(t *testing.T) {
	cleanDB(t)

	ctx := context.Background()

	id := createTestTask(t, "buy milk", task.StatusTodo)

	got, err := testRepo.GetByID(ctx, id)

	require.NoError(t, err)

	assert.Equal(t, id, got.ID)
	assert.Equal(t, "buy milk", got.Description)
	assert.Equal(t, task.StatusTodo, got.Status)
	assert.False(t, got.CreatedAt.IsZero())
	assert.False(t, got.UpdatedAt.IsZero())
}

func TestRepository_GetAll(t *testing.T) {
	cleanDB(t)

	ctx := context.Background()

	createTestTask(t, "buy milk", task.StatusTodo)
	createTestTask(t, "write code", task.StatusDone)

	got, err := testRepo.GetAll(ctx)

	require.NoError(t, err)
	require.Len(t, got, 2)

	assert.Equal(t, "buy milk", got[0].Description)
	assert.Equal(t, task.StatusTodo, got[0].Status)

	assert.Equal(t, "write code", got[1].Description)
	assert.Equal(t, task.StatusDone, got[1].Status)
}

func TestRepository_Update(t *testing.T) {
	cleanDB(t)

	ctx := context.Background()

	id := createTestTask(t, "buy milk", task.StatusTodo)

	err := testRepo.Update(ctx, task.Task{
		ID:          id,
		Description: "buy bread",
		Status:      task.StatusDone,
	})

	require.NoError(t, err)

	got, err := testRepo.GetByID(ctx, id)

	require.NoError(t, err)

	assert.Equal(t, id, got.ID)
	assert.Equal(t, "buy bread", got.Description)
	assert.Equal(t, task.StatusDone, got.Status)
}

func TestRepository_Delete(t *testing.T) {
	cleanDB(t)

	ctx := context.Background()

	id := createTestTask(t, "buy milk", task.StatusTodo)

	err := testRepo.Delete(ctx, id)

	require.NoError(t, err)

	_, err = testRepo.GetByID(ctx, id)

	assert.Error(t, err)
}

func TestRepository_StatusMapping(t *testing.T) {
	statuses := []task.Status{
		task.StatusTodo,
		task.StatusInProgress,
		task.StatusDone,
		task.StatusCanceled,
	}

	for _, status := range statuses {
		t.Run(string(status), func(t *testing.T) {
			cleanDB(t)

			id := createTestTask(t, "test task", status)

			got, err := testRepo.GetByID(
				context.Background(),
				id,
			)

			require.NoError(t, err)
			assert.Equal(t, status, got.Status)
		})
	}
}

func TestRepository_UpdateChangesUpdatedAt(t *testing.T) {
	cleanDB(t)

	ctx := context.Background()

	id := createTestTask(t, "buy milk", task.StatusTodo)

	before, err := testRepo.GetByID(ctx, id)
	require.NoError(t, err)

	time.Sleep(time.Millisecond)

	err = testRepo.Update(ctx, task.Task{
		ID:          id,
		Description: "buy bread",
		Status:      task.StatusDone,
	})

	require.NoError(t, err)

	after, err := testRepo.GetByID(ctx, id)
	require.NoError(t, err)

	assert.True(t, after.UpdatedAt.After(before.UpdatedAt))
	assert.Equal(t, before.CreatedAt, after.CreatedAt)
}
