package ws

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestWs(t *testing.T) {
	node := Node{}
	getInterfaceMethods(&node)
}

func getInterfaceMethods(i interface{}) {
	t := reflect.TypeOf(i)
	instanceValue := reflect.ValueOf(i)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		fmt.Printf("Method Name: %s\n", method.Name)
		m := instanceValue.MethodByName("CallMethod")
		m.Call([]reflect.Value{})
	}
}

type Req struct {
	VersionNumber int
}

type Response struct {
	Status string
}

type Node struct {
	Count int
}

func (n *Node) Version(req Req) (Response, error) {
	n.Count = n.Count + 1
	if req.VersionNumber > 0 {
		return Response{Status: "Success"}, nil
	}
	return Response{}, errors.New("invalid version number")
}

func TestDemo(t *testing.T) {
	// 创建Node实例
	node := &Node{}
	t.Log(node.Count)
	// 获取Node类型信息
	nodeType := reflect.TypeOf(node)
	nodeValue := reflect.ValueOf(node)
	var method reflect.Method
	var found bool
	// 遍历Node的方法，找到Version方法
	for i := 0; i < nodeType.NumMethod(); i++ {
		m := nodeType.Method(i)
		if m.Name == "Version" {
			method = m
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Version method not found")
		return
	}

	// 获取Version方法的参数类型
	if method.Type.NumIn() != 2 {
		fmt.Println("Unexpected number of input parameters")
		return
	}

	reqType := method.Type.In(1)

	// 创建Req对象并动态设置字段值
	reqValue := reflect.New(reqType).Elem()
	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)
		t.Log(field.Name)
		fieldValue := reqValue.Field(i)
		if fieldValue.CanSet() {
			switch fieldValue.Kind() {
			case reflect.Int:
				fieldValue.SetInt(1) // 动态设置为你需要的值
			case reflect.String:
				fieldValue.SetString("example")
				// 可以根据需要添加其他类型的处理
			}
		}
	}

	// 构造调用参数
	in := []reflect.Value{nodeValue, reqValue}

	// 调用方法
	out := method.Func.Call(in)

	// 处理返回值
	if len(out) != 2 {
		fmt.Println("Unexpected number of return values")
		return
	}

	// 动态创建并解析Response对象
	responseValue := out[0]
	if responseValue.Kind() == reflect.Ptr {
		responseValue = responseValue.Elem()
	}
	responseType := responseValue.Type()

	fmt.Println("Response Fields:")
	for i := 0; i < responseType.NumField(); i++ {
		field := responseType.Field(i)
		fieldValue := responseValue.Field(i)
		fmt.Printf("%s: %v\n", field.Name, fieldValue.Interface())
	}

	// 处理error
	errValue := out[1]
	if !errValue.IsNil() {
		err := errValue.Interface().(error)
		fmt.Println("Error:", err)
		return
	}
	t.Log(node.Count)
}
