package main

import (
	"log"

	micro "github.com/micro/go-micro"
	"github.com/urandom/go-micro-test/auth"
	"github.com/urandom/go-micro-test/db"
)

func main() {
	service := micro.NewService(
		micro.Name(auth.JWT),
	)

	service.Init()

	dbUser := db.NewUserClient(db.DummyDB, service.Client())

	auth.RegisterTokenHandler(service.Server(), &tokenHandler{dbUser})
	auth.RegisterUserHandler(service.Server(), &userHandler{dbUser})

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
