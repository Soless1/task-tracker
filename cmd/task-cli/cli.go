package main

import (
	"context"
	"fmt"
	"iter"
	"strconv"
	"strings"
	task2 "task-tracker/internal/domain/task"
	"task-tracker/internal/usecase/task"
)

type CLI struct {
	repo task.TaskRepository
}

func NewCLI(repo task.TaskRepository) *CLI {
	return &CLI{repo}
}

func (c *CLI) Run(ctx context.Context, args []string) {
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
	}
}

func strToID(str string) (int64, error) {
	var id uint64
	id, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}

func (c *CLI) add(ctx context.Context, args []string) {
	id, err := task.CreateTask(ctx, c.repo, strings.Join(args, " "))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Created Task id: %d\n", id)
}

func (c *CLI) update(ctx context.Context, args []string) {
	id, err := strToID(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = task.Update(ctx, c.repo, id, strings.Join(args[1:], " "))
	if err != nil {
		fmt.Println(err)
	}
}

func (c *CLI) delete(ctx context.Context, args []string) {
	id, err := strToID(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = task.Delete(ctx, c.repo, id)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *CLI) markDone(ctx context.Context, args []string) {
	id, err := strToID(args[0])

	err = task.MarkDone(ctx, c.repo, id)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *CLI) markInProgress(ctx context.Context, args []string) {
	id, err := strToID(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = task.MarkInProgress(ctx, c.repo, id)
	if err != nil {
		fmt.Println(err)
	}
}
func (c *CLI) list(ctx context.Context, args []string) {
	var tasks iter.Seq[task2.Task]
	var err error
	if len(args) == 0 {
		tasks, err = task.List(ctx, c.repo, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		var status task2.Status
		switch args[0] {
		case "todo":
			status = task2.StatusTodo
		case "done":
			status = task2.StatusDone
		case "in-progress":
			status = task2.StatusInProgress
		default:
			fmt.Println("list error: Invalid status: " + args[0])
			return
		}
		tasks, err = task.List(ctx, c.repo, &status)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Printf("%-4s %-15s %-30s %s\n",
		"ID", "STATUS", "DESCRIPTION", "CREATED",
	)
	for t := range tasks {
		fmt.Printf("%-4d %-15s %-30s %s\n",
			t.ID,
			t.Status,
			t.Description,
			t.CreatedAt.Format("2006-01-02 15:04"),
		)
	}
}
