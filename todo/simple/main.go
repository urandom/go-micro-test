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

	token := auth.NewTokenClient(auth.JWT, service.Client())
	dbTodo := db.NewTodoClient(db.DummyDB, service.Client())

	todo.RegisterTodoHandler(service.Server(), &todoHandler{dbTodo, token})

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
