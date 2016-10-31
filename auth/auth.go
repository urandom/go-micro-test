package auth

//go:generate protoc --go_out=plugins=micro:. auth.proto

const (
	JWT = "org.sugr.micro.test.auth.jwt"
)

var (
	SecretKey = []byte("secret-key")
)
