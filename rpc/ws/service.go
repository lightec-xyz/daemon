package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"unicode"
)

// Service todo
// just a simple rpc service,need more check and test
type Service struct {
	calls *sync.Map // map[string]*call
}

func NewService(h interface{}) *Service {
	calls := new(sync.Map)
	receiver := reflect.ValueOf(h)
	typ := receiver.Type()
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		if method.PkgPath != "" {
			continue
		}
		name := formatName(method.Name)
		cl := newCall(name, receiver, method.Func)
		if cl == nil {
			continue
		}
		calls.Store(name, cl)
	}
	return &Service{
		calls: calls,
	}
}

func (h *Service) Call(name string, argsByte []byte) (interface{}, error) {
	call, ok := h.GetCall(name)
	if !ok {
		return nil, fmt.Errorf("no such method: %s", name)
	}
	result, err := call.call(argsByte)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (h *Service) GetCall(method string) (*call, bool) {
	value, ok := h.calls.Load(method)
	if !ok {
		return nil, false
	}
	call, ok := value.(*call)
	if !ok {
		return nil, false
	}
	return call, true

}

func (h *Service) Check(method string) bool {
	_, ok := h.calls.Load(method)
	return ok
}

type call struct {
	name     string
	fn       reflect.Value  // the function
	rcvr     reflect.Value  // receiver object of method, set if fn is method
	argTypes []reflect.Type // input argument types
}

func (c *call) call(param []byte) (interface{}, error) {
	// todo
	var fullArgs []reflect.Value
	if c.rcvr.IsValid() {
		fullArgs = append(fullArgs, c.rcvr)
	}
	if len(param) != 0 {
		args, err := parseJsonToArgs(param, c.argTypes)
		if err != nil {
			return nil, err
		}
		fullArgs = append(fullArgs, args...)
	}
	out := c.fn.Call(fullArgs)
	if len(out) == 0 {
		return nil, nil
	} else if len(out) == 1 {
		return out[0].Interface(), nil
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

func parseJsonToArgs(params []byte, types []reflect.Type) ([]reflect.Value, error) {
	// todo
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
