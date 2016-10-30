package main

import (
	"log"

	micro "github.com/micro/go-micro"
	"github.com/urandom/go-micro-test/db"
)

func main() {
	service := micro.NewService(
		micro.Name(db.DummyDB),
		micro.Version("0.0.1"),
	)

	service.Init()

	db.RegisterTodoHandler(service.Server(), &todoHandler{})
	db.RegisterProfilerHandler(service.Server(), &profilerHandler{})

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
