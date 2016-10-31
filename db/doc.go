package db

//go:generate protoc --go_out=plugins=micro:. db.proto

const (
	DummyDB = "org.sugr.micro.test.db.dummy"
)
