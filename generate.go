package main

//go:generate sqlc generate
//go:generate mockgen -source=repository.go -destination=mocks/task_repository.go -package=mocks
