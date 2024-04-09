package common

import (
	"container/list"
	"fmt"
	"testing"
)

func TestUuid(t *testing.T) {
	uuid, err := Uuid()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(uuid)
}

func TestDemo(t *testing.T) {
	// 创建一个列表
	myList := list.New()

	// 向列表中添加一些元素
	myList.PushBack("Hello")
	myList.PushBack("Go")
	myList.PushBack("World")

	// 遍历列表并打印每个元素
	for element := myList.Front(); element != nil; element = element.Next() {
		fmt.Println(element.Value)
	}
}
