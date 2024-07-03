package ws

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestHandler(t *testing.T) {
	node := testHandler{}
	handler := newHandler(&node)

	info, err := handler.call("info")
	if err != nil {
		t.Error(err)
	}
	t.Log(info)

}

func TestDemo(t *testing.T) {
	data := []byte(`{"name":"able","age":10}`)
	//data := []byte(`[1,"age",3,4,5,6,"able"]`)
	decoder := json.NewDecoder(bytes.NewReader(data))
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			t.Error(err)
		}
		t.Log(token)
	}

}
