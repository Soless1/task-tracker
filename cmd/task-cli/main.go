package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"task-tracker/internal/repository/json"
	"task-tracker/internal/repository/postgres"
	"task-tracker/internal/usecase/task"
	"task-tracker/migrations"

	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "task-tracker/docs"
	httptransport "task-tracker/internal/transport/http"
)

// @title Task Tracker API
// @version 1.0
// @description REST API for Task Tracker
// @host localhost:8080
// @BasePath /
func main() {
	ctx := context.Background()

	mode := flag.String(
		"mode",
		"web",
		"application mode: cli or web",
	)

	storage := flag.String(
		"storage",
		"postgres",
		"storage: json or postgres",
	)

	flag.Parse()

	repo, cleanup, err := createRepository(ctx, *storage)
	if err != nil {
		fmt.Println("failed to create repository:", err)
		return
	}
	defer cleanup()

	switch *mode {
	case "cli":
		runCLI(ctx, repo, flag.Args())

	case "web":
		runWebServer(repo)

	default:
		fmt.Printf("unknown mode: %q\n", *mode)
		fmt.Println("available modes: cli, web")
	}
}

func createRepository(
	ctx context.Context,
	storage string,
) (task.TaskRepository, func(), error) {

	switch storage {
	case "json":
		repo, err := json.NewJSONTaskRepository()
		if err != nil {
			return nil, func() {}, err
		}

		return repo, func() {}, nil

	case "postgres":
		pool, err := pgxpool.New(
			ctx,
			"postgres://postgres:postgres@localhost:5432/task_tracker?sslmode=disable",
		)
		if err != nil {
			return nil, func() {}, err
		}

		if err := migrations.Run(pool); err != nil {
			pool.Close()
			return nil, func() {}, err
		}

		repo := postgres.NewRepository(pool)

		return repo, func() {
			pool.Close()
		}, nil

	default:
		return nil, func() {}, fmt.Errorf(
			"unknown storage: %q (available: json, postgres)",
			storage,
		)
	}
}

func runCLI(
	ctx context.Context,
	repo task.TaskRepository,
	args []string,
) {
	cli := NewCLI(repo)
	cli.Run(ctx, args)
}

func runWebServer(repo task.TaskRepository) {
	h := httptransport.NewHandler(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/tasks", h.GetTasks)
	mux.HandleFunc("POST /api/tasks", h.CreateTask)
	mux.HandleFunc("GET /api/tasks/{id}", h.GetTask)
	mux.HandleFunc("PUT /api/tasks/{id}", h.UpdateTask)
	mux.HandleFunc("DELETE /api/tasks/{id}", h.DeleteTask)
	mux.HandleFunc("POST /api/tasks/{id}/done", h.MarkDone)
	mux.HandleFunc("POST /api/tasks/{id}/in-progress", h.MarkInProgress)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux.Handle("/", http.FileServer(http.Dir("./web")))

	const addr = ":8080"

	fmt.Println("server started on", addr)
	fmt.Println("web:     http://localhost:8080")
	fmt.Println("swagger: http://localhost:8080/swagger/index.html")

	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Println("server stopped:", err)
	}
}
