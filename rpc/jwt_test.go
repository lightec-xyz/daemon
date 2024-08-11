package rpc

import (
	"testing"
)

// var secretKey = []byte("your-secret-key")
var secretKey = []byte("f9962336ca15bdd2acd61edfc6857fe733ef36a3c1380acf5f91c17347df93e5")

func TestNewJwt(t *testing.T) {
	jwt, err := CreateJWT(secretKey, JwtPermission)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(jwt)
	claims, err := VerifyJWT(secretKey, jwt)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(claims)
}
