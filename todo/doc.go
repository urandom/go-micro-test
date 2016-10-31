package todo

//go:generate protoc --go_out=plugins=micro:. todo.proto

const (
	Simple = "org.sugr.micro.test.todo.simple"
)
