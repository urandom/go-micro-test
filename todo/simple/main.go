package main

import (
	"log"

	micro "github.com/micro/go-micro"
	"github.com/urandom/go-micro-test/auth"
	"github.com/urandom/go-micro-test/db"
	"github.com/urandom/go-micro-test/todo"
)

func main() {
	service := micro.NewService(
		micro.Name(todo.Simple),
	)

	service.Init()

	checker := auth.NewCheckerClient(auth.JWT, service.Client())
	dbTodo := db.NewTodoClient(db.DummyDB, service.Client())

	todo.RegisterServiceHandler(service.Server(), &todoHandler{dbTodo, checker})

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
