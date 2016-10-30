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

	dbProfiler := db.NewProfilerClient(db.DummyDB, service.Client())

	auth.RegisterGeneratorHandler(service.Server(), &generatorHandler{dbProfiler})
	auth.RegisterCheckerHandler(service.Server(), &checkerHandler{dbProfiler})
	auth.RegisterProfilerHandler(service.Server(), &profilerHandler{dbProfiler})

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
