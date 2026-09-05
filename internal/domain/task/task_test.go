package task

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	createdAt := time.Now()

	got := CreateTask("buy milk", createdAt, 1)

	assert.Equal(t, int64(1), got.ID)
	assert.Equal(t, "buy milk", got.Description)
	assert.Equal(t, StatusTodo, got.Status)
	assert.Equal(t, createdAt, got.CreatedAt)
	assert.Equal(t, createdAt, got.UpdatedAt)
}

func TestTask_MarkDone(t *testing.T) {
	createdAt := time.Now()
	tk := CreateTask("buy milk", createdAt, 1)

	before := time.Now()
	tk.MarkDone()
	after := time.Now()

	assert.Equal(t, StatusDone, tk.Status)
	assert.False(t, tk.UpdatedAt.Before(before))
	assert.False(t, tk.UpdatedAt.After(after))
}

func TestTask_MarkInProgress(t *testing.T) {
	tk := CreateTask("buy milk", time.Now(), 1)

	tk.MarkInProgress()

	assert.Equal(t, StatusInProgress, tk.Status)
}

func TestTask_MarkCanceled(t *testing.T) {
	tk := CreateTask("buy milk", time.Now(), 1)

	tk.MarkCanceled()

	assert.Equal(t, StatusCanceled, tk.Status)
}

func TestTask_ChangeDescription(t *testing.T) {
	tk := CreateTask("buy milk", time.Now(), 1)

	tk.ChangeDescription("wash dishes")

	assert.Equal(t, "wash dishes", tk.Description)
}
