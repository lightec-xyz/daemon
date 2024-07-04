package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"unicode"
)

/*
todo
just a simple rpc service,need more check and test
*/
type Service struct {
	calls map[string]*call
}

func NewService(h interface{}) *Service {
	calls := make(map[string]*call)
	receiver := reflect.ValueOf(h)
	typ := receiver.Type()
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		if method.PkgPath != "" {
			continue
		}
		name := formatName(method.Name)
		cb := newCall(name, receiver, method.Func)
		if cb == nil {
			continue
		}
		calls[name] = cb
	}

	return &Service{
		calls: calls,
	}
}

func (h *Service) Call(name string, args ...interface{}) (interface{}, error) {
	var argsByte []byte
	var err error
	if len(args) != 0 {
		argsByte, err = json.Marshal(args)
		if err != nil {
			return nil, err
		}
	}
	call, ok := h.calls[name]
	if !ok {
		return nil, fmt.Errorf("no such method: %s", name)
	}
	result, err := call.call(argsByte)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type call struct {
	name     string
	fn       reflect.Value  // the function
	rcvr     reflect.Value  // receiver object of method, set if fn is method
	argTypes []reflect.Type // input argument types
}

func (c *call) call(obj []byte) (interface{}, error) {
	var fullArgs []reflect.Value
	if c.rcvr.IsValid() {
		fullArgs = append(fullArgs, c.rcvr)
	}
	if len(obj) != 0 {
		args, err := parseJsonArrayToArgs(obj, c.argTypes)
		if err != nil {
			return nil, err
		}
		fullArgs = append(fullArgs, args...)
	}
	out := c.fn.Call(fullArgs)
	if len(out) == 0 {
		return nil, nil
	} else if len(out) == 1 {
		return nil, out[0].Interface().(error)
	} else if len(out) == 2 {
		if out[1].IsNil() {
			return out[0].Interface(), nil
		}
		return nil, out[1].Interface().(error)
	} else {
		return nil, fmt.Errorf("return value more than 2")
	}
}
func (c *call) makeArgTypes() {
	fntype := c.fn.Type()
	// Skip receiver and context.Context parameter (if present).
	firstArg := 0
	if c.rcvr.IsValid() {
		firstArg++
	}
	// Add all remaining parameters.
	c.argTypes = make([]reflect.Type, fntype.NumIn()-firstArg)
	for i := firstArg; i < fntype.NumIn(); i++ {
		c.argTypes[i-firstArg] = fntype.In(i)
	}
}

func newCall(name string, receiver, fn reflect.Value) *call {
	c := &call{name: name, fn: fn, rcvr: receiver}
	c.makeArgTypes()
	return c
}

func formatName(name string) string {
	ret := []rune(name)
	if len(ret) > 0 {
		ret[0] = unicode.ToLower(ret[0])
	}
	return string(ret)
}

func parseJsonArrayToArgs(params []byte, types []reflect.Type) ([]reflect.Value, error) {
	dec := json.NewDecoder(bytes.NewReader(params))
	args := make([]reflect.Value, 0, len(types))
	token, err := dec.Token()
	if err != nil {
		return args, err
	}
	if token != json.Delim('[') {
		return args, fmt.Errorf("expected '['")
	}
	for i := 0; dec.More(); i++ {
		if i >= len(types) {
			return args, fmt.Errorf("too many arguments, want at most %d", len(types))
		}
		argval := reflect.New(types[i])
		if err := dec.Decode(argval.Interface()); err != nil {
			return args, fmt.Errorf("invalid argument %d: %v", i, err)
		}
		if argval.IsNil() && types[i].Kind() != reflect.Ptr {
			return args, fmt.Errorf("missing value for required argument %d", i)
		}
		args = append(args, argval.Elem())
	}
	// Read end of args array.
	_, err = dec.Token()
	return args, err
}
