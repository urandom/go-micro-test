package main

import (
	"log"
	"net/http"

	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
	"github.com/urandom/go-micro-test/auth"
	"github.com/urandom/go-micro-test/db"
	"github.com/urandom/go-micro-test/todo"
)

func main() {
	service := micro.NewService(
		micro.Name(db.DummyDB),
		micro.Flags(
			cli.StringFlag{
				Name:  "listen-on",
				Usage: "Start on http server on this address",
				Value: ":8080",
			},
		),
	)

	service.Init(micro.Action(func(c *cli.Context) {
		token := auth.NewTokenClient(auth.JWT, service.Client())
		user := auth.NewUserClient(auth.JWT, service.Client())
		todo := todo.NewTodoClient(todo.Simple, service.Client())

		http.Handle("/user", profileHandler{user})
		http.Handle("/user/login", loginHandler{token})
		http.Handle("/user/todo", todoHandler{todo})

		err := http.ListenAndServe(c.String("listen-on"), nil)
		if err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}))
}
