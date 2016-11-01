// Package auth provides the various service clients for working with
// authentication.
package auth

//go:generate protoc --go_out=plugins=micro:. auth.proto

const (
	// JWT is the name of the auth service that works with jwt tokens.
	JWT = "org.sugr.micro.test.auth.jwt"
)

var (
	// SecretKey is the secret key used by auth service for encrypting their
	// tokens.
	SecretKey = []byte("secret-key")
)
