package main

import (
	"log"
	"net/http"

	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
	"github.com/urandom/go-micro-test/auth"
	"github.com/urandom/go-micro-test/db"
	"github.com/urandom/go-micro-test/todo"
	"golang.org/x/net/context"
)

func main() {
	service := micro.NewService(
		micro.Name(db.DummyDB),
		micro.Flags(
			cli.BoolFlag{
				Name:  "generate-jwt",
				Usage: "Generate a jwt token for a user and pass arguments",
			},
			cli.BoolFlag{
				Name:  "check-jwt",
				Usage: "Check a supplied jwt token",
			},
			cli.BoolFlag{
				Name:  "todo",
				Usage: "Fetches all todo items given a jwt token",
			},
			cli.BoolFlag{
				Name:  "profile",
				Usage: "Fetches the profile data given a jwt token",
			},
			cli.StringFlag{
				Name:  "listen-on",
				Usage: "Start on http server on this address",
				Value: ":8080",
			},
		),
	)

	service.Init(micro.Action(func(c *cli.Context) {
		if c.Bool("generate-jwt") {
			if c.NArg() < 2 {
				log.Fatalln("Insufficient number of arguments")
			}
			generator := auth.NewGeneratorClient(auth.JWT, service.Client())

			user := c.Args()[0]
			pass := c.Args()[1]
			resp, err := generator.Generate(context.TODO(), &auth.GenerateRequest{user, pass})
			if err != nil {
				log.Fatalf("Generating jwt token for %s:%s: %+v", user, pass, err)
			}

			if resp.Exists {
				log.Printf("JWT token: %s", resp.Auth)
			} else {
				log.Printf("User %s doesn't exist", user)
			}
		} else if c.Bool("check-jwt") {
			if c.NArg() < 1 {
				log.Fatalln("Insufficient number of arguments")
			}
			checker := auth.NewCheckerClient(auth.JWT, service.Client())

			token := c.Args()[0]
			resp, err := checker.Check(context.TODO(), &auth.CheckRequest{token})
			if err != nil {
				log.Fatalf("Checking jwt token %s: %+v", token, err)
			}

			if resp.Valid {
				log.Println("JWT token is valid")
			} else if resp.Expired {
				log.Println("JWT token has expired")
			} else {
				log.Println("Token is not valid")
			}
		} else if c.Bool("todo") {
			if c.NArg() < 1 {
				log.Fatalln("Insufficient number of arguments")
			}

			s := todo.NewServiceClient(todo.Simple, service.Client())
			token := c.Args()[0]
			resp, err := s.List(context.TODO(), &todo.Request{token, 0, 0})
			if err != nil {
				log.Fatalf("Getting todo items for %s: %+v", token, err)
			}

			if resp.Valid {
				log.Printf("Todo item count: %d", len(resp.Items))
				for _, e := range resp.Items {
					log.Printf("Todo item %s: %s", e.Title, e.Body)
				}
			} else {
				log.Println("Token is not valid")
			}
		} else if c.Bool("profile") {
			if c.NArg() < 1 {
				log.Fatalln("Insufficient number of arguments")
			}
			profiler := auth.NewProfilerClient(auth.JWT, service.Client())

			token := c.Args()[0]
			resp, err := profiler.UserProfile(context.TODO(), &auth.ProfileRequest{token})
			if err != nil {
				log.Fatalf("Getting user profile for %s: %+v", token, err)
			}

			if resp.Exists {
				log.Printf("Profile: %#v", resp.Profile)
			} else {
				log.Println("User doesn't exist")
			}
		} else {
			generator := auth.NewGeneratorClient(auth.JWT, service.Client())
			profiler := auth.NewProfilerClient(auth.JWT, service.Client())
			todo := todo.NewServiceClient(todo.Simple, service.Client())

			http.Handle("/auth", loginHandler{generator})
			http.Handle("/user", profileHandler{profiler})
			http.Handle("/user/todo", todoHandler{todo})

			http.ListenAndServe(c.String("listen-on"), nil)
		}
	}))
}
