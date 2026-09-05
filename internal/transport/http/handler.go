package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	domainTask "task-tracker/internal/domain/task"
	taskUsecase "task-tracker/internal/usecase/task"
)

type CreateTaskRequest struct {
	Description string `json:"description" example:"Learn REST"`
}

type CreateTaskResponse struct {
	ID int64 `json:"id" example:"1"`
}

type UpdateTaskRequest struct {
	Description string `json:"description" example:"Learn REST and Swagger"`
}
type Handler struct {
	repo taskUsecase.TaskRepository
}

func NewHandler(repo taskUsecase.TaskRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func parseStatus(value string) (*domainTask.Status, error) {
	if value == "" {
		return nil, nil
	}

	status := domainTask.Status(value)

	switch status {
	case domainTask.StatusTodo,
		domainTask.StatusInProgress,
		domainTask.StatusDone,
		domainTask.StatusCanceled:
		return &status, nil
	default:
		return nil, fmt.Errorf("unknown status: %s", value)
	}
}

// GetTasks godoc
// @Summary Get all tasks
// @Description Returns all tasks
// @Tags tasks
// @Produce json
// @Success 200 {array} task.Task
// @Failure 500 {string} string
// @Router /api/tasks [get]
func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	status, err := parseStatus(r.URL.Query().Get("status"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	seq, err := taskUsecase.List(r.Context(), h.repo, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tasks := make([]domainTask.Task, 0)

	seq(func(t domainTask.Task) bool {
		tasks = append(tasks, t)
		return true
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// CreateTask godoc
// @Summary Create a task
// @Description Creates a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body CreateTaskRequest true "Task data"
// @Success 201 {object} CreateTaskResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /api/tasks [post]
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := taskUsecase.CreateTask(
		r.Context(),
		h.repo,
		request.Description,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(CreateTaskResponse{
		ID: id,
	})
}

// GetTask godoc
// @Summary Get task by ID
// @Description Returns a task by its ID
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} task.Task
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /api/tasks/{id} [get]
func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	t, err := taskUsecase.GetTask(r.Context(), h.repo, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(t)
}

// UpdateTask godoc
// @Summary Update a task
// @Description Updates a task description
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param request body UpdateTaskRequest true "Task data"
// @Success 204
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/tasks/{id} [put]
func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	var request UpdateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = taskUsecase.Update(
		r.Context(),
		h.repo,
		id,
		request.Description,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Deletes a task by its ID
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 204
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /api/tasks/{id} [delete]
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	err = taskUsecase.Delete(r.Context(), h.repo, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// MarkDone godoc
// @Summary Mark task as done
// @Description Marks a task as completed
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 204
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /api/tasks/{id}/done [post]
func (h *Handler) MarkDone(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	err = taskUsecase.MarkDone(r.Context(), h.repo, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// MarkInProgress godoc
// @Summary Mark task as in progress
// @Description Marks a task as in progress
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 204
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /api/tasks/{id}/in-progress [post]
func (h *Handler) MarkInProgress(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	err = taskUsecase.MarkInProgress(r.Context(), h.repo, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
