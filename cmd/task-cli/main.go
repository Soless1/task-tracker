package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"task-tracker/internal/repository/postgres"
	"task-tracker/migrations"

	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "task-tracker/docs"
	httptransport "task-tracker/internal/transport/http"
)

//func main() {
//	var repo task.TaskRepository
//	repo, err := json.NewJSONTaskRepository()
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	cli := NewCLI(repo)
//	cli.Run(os.Args[1:])
//}

// @title Task Tracker API
// @version 1.0
// @description REST API for Task Tracker
// @host localhost:8080
// @BasePath /
func main() {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/task_tracker?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pool.Close()

	err = migrations.Run(pool)
	if err != nil {
		fmt.Println(err)
		return
	}

	repo := postgres.NewRepository(pool)

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

	fmt.Println("server started on :8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err)
	}

	cli := NewCLI(repo)
	cli.Run(ctx, os.Args[1:])
}
