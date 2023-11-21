package rpc

import "testing"

var err error
var client *Client

func init() {
	client, err = NewClient("http://127.0.0.1:8445")
	if err != nil {
		panic(err)
	}
}
func TestClient_HelloWorld(t *testing.T) {
	var result string
	name := "red"
	age := 100
	err := client.Call(&result, "zkbtc_helloWorld", &name, &age)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
