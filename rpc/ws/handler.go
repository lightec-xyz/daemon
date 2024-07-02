package ws

import (
	"fmt"
	"reflect"
)

// todo

type handler struct {
	calls map[string]*call
}

func newHandler(h interface{}) *handler {
	calls := make(map[string]*call)
	rTypes := reflect.TypeOf(h)
	rValues := reflect.ValueOf(h)
	for i := 0; i < rTypes.NumMethod(); i++ {
		method := rTypes.Method(i)
		calls[method.Name] = &call{
			name:   method.Name,
			method: method,
			root:   rValues,
		}
	}
	return &handler{
		calls: calls,
	}
}

func (h *handler) call(name string) interface{} {
	return h.calls[name].call(nil)
}

type call struct {
	name   string
	args   []reflect.Value
	method reflect.Method
	root   reflect.Value
}

func (c *call) call(obj interface{}) interface{} {
	reqType := c.method.Type.In(1)
	reqValue := reflect.New(reqType).Elem()
	for i := 0; i < reqType.NumField(); i++ {
		//field := reqType.Field(i)
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
	out := c.method.Func.Call([]reflect.Value{c.root, reqValue})
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

	}
	return out[0].Interface()
}
