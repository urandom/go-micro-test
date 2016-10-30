package auth

//go:generate protoc --go_out=plugins=micro:. auth.proto

const (
	JWT = "micro-test-auth-jwt"
)

var (
	SecretKey = []byte("secret-key")
)
