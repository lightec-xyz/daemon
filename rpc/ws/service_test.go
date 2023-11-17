package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestService(t *testing.T) {
	node := testService{}
	service := NewService(&node)
	info, err := service.Call("info", nil)
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
	res, err := service.Call("add", getArg("red", 100))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
	result, err := service.Call("version", getArg(Request{Height: 10001}))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)

}

func TestDemo(t *testing.T) {
	//data := []byte(`{"name":"able","age":10}`)
	//data := []byte(`[1,"age",3,4,5,6,"able"]`)
	data := []byte(`[{"Height":100}]`)
	decoder := json.NewDecoder(bytes.NewReader(data))
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			t.Error(err)
		}
		t.Log(token)
	}

}

func getArg(args ...interface{}) []byte {
	marshal, _ := json.Marshal(args)
	return marshal
}

type testService struct {
	name  string
	count int
}

func (t *testService) Version(req Request) (*Response, error) {
	return &Response{Msg: fmt.Sprintf("version height: %v", req.Height)}, nil
}

func (t *testService) Add(name string, count int) (string, error) {
	t.name = name
	t.count = count
	res := fmt.Sprintf("name: %v,count: %v", name, count)
	return res, nil
}

func (t *testService) Info() (string, error) {
	info := fmt.Sprintf("info test %v:%v", t.name, t.count)
	return info, nil
}

func (t *testService) Test() error {
	fmt.Printf("test:")
	return nil
}

type Request struct {
	Height int
}

type Response struct {
	Msg string
}
