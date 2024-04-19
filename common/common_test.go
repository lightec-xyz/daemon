package common

import (
	"strings"
	"testing"
)

func TestUuid(t *testing.T) {
	uuid, err := Uuid()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(uuid)
}

func TestDemo001(t *testing.T) {
	fold := strings.EqualFold("genesis", "genesis")
	t.Log(fold)
}

func TestDemo(t *testing.T) {
	t.Log(GetNearTxSlot(63))
	t.Log(GetNearTxSlot(64))
	t.Log(GetNearTxSlot(65))

}
