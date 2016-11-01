// Package db provides the various service clients for getting data from an
// implementation.
package db

//go:generate protoc --go_out=plugins=micro:. db.proto

const (
	// DummyDB is the name of the dummy service implementation of the
	// various handlers.
	DummyDB = "org.sugr.micro.test.db.dummy"
)
