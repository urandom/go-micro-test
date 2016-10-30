package db

//go:generate protoc --go_out=plugins=micro:. db.proto

const (
	DummyDB = "micro-test-dummy-db"
)
