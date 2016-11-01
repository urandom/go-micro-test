// Package todo provides a service client for obtaining a user's todo list
// items.
package todo

//go:generate protoc --go_out=plugins=micro:. todo.proto

const (
	// Simple is the name of the service implementation of the TodoHandler
	Simple = "org.sugr.micro.test.todo.simple"
)
