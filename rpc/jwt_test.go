package rpc

import (
	"encoding/hex"
	"testing"
	"time"
)

// var secretKey = []byte("your-secret-key")
var secretKey = ""

func TestNewJwt(t *testing.T) {
	hexSec, err := hex.DecodeString(secretKey)
	if err != nil {
		t.Fatal(err)
	}
	jwt, err := CreateJWT(hexSec, JwtPermission, 100*24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(jwt)
	claims, err := VerifyJWT(hexSec, jwt)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(claims)
}
