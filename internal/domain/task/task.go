package task

import "time"

type Task struct {
	ID          int64
	Description string
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *Task) MarkDone() {
	t.Status = StatusDone
	t.UpdatedAt = time.Now()
}

func (t *Task) MarkCanceled() {
	t.Status = StatusCanceled
	t.UpdatedAt = time.Now()
}

func (t *Task) MarkInProgress() {
	t.Status = StatusInProgress
	t.UpdatedAt = time.Now()
}
func CreateTask(description string, createdAt time.Time, ID int64) Task {
	return Task{
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
		ID:          ID,
	}
}

func (t *Task) ChangeDescription(description string) {
	t.Description = description
}
