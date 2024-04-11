package common

import (
	"container/list"
	"fmt"
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
	myList := list.New()
	myList.PushBack("Hello")
	myList.PushBack("Go")
	myList.PushBack("World")
	for element := myList.Front(); element != nil; element = element.Next() {
		fmt.Println(element.Value)
	}
}
