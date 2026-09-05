package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	taskdomain "task-tracker/internal/domain/task"
	taskusecase "task-tracker/internal/usecase/task"
)

type CLI struct {
	repo taskusecase.TaskRepository
}

func NewCLI(repo taskusecase.TaskRepository) *CLI {
	return &CLI{
		repo: repo,
	}
}

func (c *CLI) Run(ctx context.Context, args []string) {
	if len(args) == 0 {
		c.printHelp()
		return
	}

	switch args[0] {
	case "add":
		c.add(ctx, args[1:])

	case "update":
		c.update(ctx, args[1:])

	case "delete":
		c.delete(ctx, args[1:])

	case "mark-done":
		c.markDone(ctx, args[1:])

	case "mark-in-progress":
		c.markInProgress(ctx, args[1:])

	case "list":
		c.list(ctx, args[1:])

	default:
		fmt.Printf("unknown command: %s\n", args[0])
		c.printHelp()
	}
}

func strToID(str string) (int64, error) {
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid task ID %q: %w", str, err)
	}

	if id <= 0 {
		return 0, fmt.Errorf("task ID must be positive")
	}

	return id, nil
}

func (c *CLI) add(ctx context.Context, args []string) {
	if len(args) == 0 {
		fmt.Println("usage: add <description>")
		return
	}

	description := strings.Join(args, " ")

	id, err := taskusecase.CreateTask(ctx, c.repo, description)
	if err != nil {
		fmt.Println("add error:", err)
		return
	}

	fmt.Printf("Created task with ID: %d\n", id)
}

func (c *CLI) update(ctx context.Context, args []string) {
	if len(args) < 2 {
		fmt.Println("usage: update <id> <description>")
		return
	}

	id, err := strToID(args[0])
	if err != nil {
		fmt.Println("update error:", err)
		return
	}

	description := strings.Join(args[1:], " ")

	err = taskusecase.Update(ctx, c.repo, id, description)
	if err != nil {
		fmt.Println("update error:", err)
		return
	}

	fmt.Printf("Updated task %d\n", id)
}

func (c *CLI) delete(ctx context.Context, args []string) {
	if len(args) != 1 {
		fmt.Println("usage: delete <id>")
		return
	}

	id, err := strToID(args[0])
	if err != nil {
		fmt.Println("delete error:", err)
		return
	}

	err = taskusecase.Delete(ctx, c.repo, id)
	if err != nil {
		fmt.Println("delete error:", err)
		return
	}

	fmt.Printf("Deleted task %d\n", id)
}

func (c *CLI) markDone(ctx context.Context, args []string) {
	if len(args) != 1 {
		fmt.Println("usage: mark-done <id>")
		return
	}

	id, err := strToID(args[0])
	if err != nil {
		fmt.Println("mark-done error:", err)
		return
	}

	err = taskusecase.MarkDone(ctx, c.repo, id)
	if err != nil {
		fmt.Println("mark-done error:", err)
		return
	}

	fmt.Printf("Task %d marked as done\n", id)
}

func (c *CLI) markInProgress(ctx context.Context, args []string) {
	if len(args) != 1 {
		fmt.Println("usage: mark-in-progress <id>")
		return
	}

	id, err := strToID(args[0])
	if err != nil {
		fmt.Println("mark-in-progress error:", err)
		return
	}

	err = taskusecase.MarkInProgress(ctx, c.repo, id)
	if err != nil {
		fmt.Println("mark-in-progress error:", err)
		return
	}

	fmt.Printf("Task %d marked as in progress\n", id)
}

func (c *CLI) list(ctx context.Context, args []string) {
	var status *taskdomain.Status

	if len(args) > 1 {
		fmt.Println("usage: list [status]")
		return
	}

	if len(args) == 1 {
		parsedStatus, err := parseStatus(args[0])
		if err != nil {
			fmt.Println("list error:", err)
			return
		}

		status = parsedStatus
	}

	tasks, err := taskusecase.List(ctx, c.repo, status)
	if err != nil {
		fmt.Println("list error:", err)
		return
	}

	fmt.Printf(
		"%-4s %-15s %-30s %s\n",
		"ID",
		"STATUS",
		"DESCRIPTION",
		"CREATED",
	)

	for t := range tasks {
		fmt.Printf(
			"%-4d %-15s %-30s %s\n",
			t.ID,
			t.Status,
			t.Description,
			t.CreatedAt.Format("2006-01-02 15:04"),
		)
	}
}

func parseStatus(value string) (*taskdomain.Status, error) {
	status := taskdomain.Status(value)

	switch status {
	case taskdomain.StatusTodo,
		taskdomain.StatusInProgress,
		taskdomain.StatusDone:
		return &status, nil

	default:
		return nil, fmt.Errorf("invalid status: %s", value)
	}
}

func (c *CLI) printHelp() {
	fmt.Println("Task Tracker CLI")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  add <description>")
	fmt.Println("  update <id> <description>")
	fmt.Println("  delete <id>")
	fmt.Println("  mark-done <id>")
	fmt.Println("  mark-in-progress <id>")
	fmt.Println("  list [status]")
	fmt.Println()
	fmt.Println("Statuses:")
	fmt.Println("  todo")
	fmt.Println("  in-progress")
	fmt.Println("  done")
}
