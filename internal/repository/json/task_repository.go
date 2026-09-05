package json

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"task-tracker/internal/domain/task"
)

const JSONPath = "task.json"

var ErrTaskNotFound = errors.New("task not found")

type JSONTaskRepository struct {
	path string
	ID   int64
}

func NewJSONTaskRepository() (*JSONTaskRepository, error) {
	return NewJSONTaskRepositoryWithPath(JSONPath)
}

func NewJSONTaskRepositoryWithPath(path string) (*JSONTaskRepository, error) {
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		if err := os.WriteFile(path, []byte("{}"), 0644); err != nil {
			return nil, err
		}
	}

	tasks, err := loadTasks(path)
	if err != nil {
		return nil, err
	}

	return &JSONTaskRepository{
		path: path,
		ID:   getNextID(tasks),
	}, nil
}

func getNextID(tasks map[string]task.Task) int64 {
	var maxID int64

	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}

	return maxID + 1
}

func (r *JSONTaskRepository) NextID(ctx context.Context) int64 {
	id := r.ID
	r.ID++

	return id
}

func (r *JSONTaskRepository) Create(
	ctx context.Context,
	t task.Task,
) (int64, error) {
	tasks, err := r.loadTasks()
	if err != nil {
		return 0, err
	}

	tasks[idToString(t.ID)] = t

	if err := r.saveTasks(tasks); err != nil {
		return 0, err
	}

	return t.ID, nil
}

func (r *JSONTaskRepository) GetByID(
	ctx context.Context,
	id int64,
) (task.Task, error) {
	tasks, err := r.loadTasks()
	if err != nil {
		return task.Task{}, err
	}

	t, ok := tasks[idToString(id)]
	if !ok {
		return task.Task{}, ErrTaskNotFound
	}

	return t, nil
}

func (r *JSONTaskRepository) GetAll(
	ctx context.Context,
) ([]task.Task, error) {
	tasks, err := r.loadTasks()
	if err != nil {
		return nil, err
	}

	result := make([]task.Task, 0, len(tasks))

	for _, t := range tasks {
		result = append(result, t)
	}

	return result, nil
}

func (r *JSONTaskRepository) Update(
	ctx context.Context,
	t task.Task,
) error {
	tasks, err := r.loadTasks()
	if err != nil {
		return err
	}

	key := idToString(t.ID)

	if _, ok := tasks[key]; !ok {
		return ErrTaskNotFound
	}

	tasks[key] = t

	return r.saveTasks(tasks)
}

func (r *JSONTaskRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	tasks, err := r.loadTasks()
	if err != nil {
		return err
	}

	key := idToString(id)

	if _, ok := tasks[key]; !ok {
		return ErrTaskNotFound
	}

	delete(tasks, key)

	return r.saveTasks(tasks)
}

func (r *JSONTaskRepository) loadTasks() (map[string]task.Task, error) {
	return loadTasks(r.path)
}

func (r *JSONTaskRepository) saveTasks(tasks map[string]task.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.path, data, 0644)
}

func loadTasks(path string) (map[string]task.Task, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var tasks map[string]task.Task

	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = make(map[string]task.Task)
	}

	return tasks, nil
}

func idToString(id int64) string {
	return strconv.FormatInt(id, 10)
}
