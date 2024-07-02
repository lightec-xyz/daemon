package ws

import (
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
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		fmt.Printf("Method Name: %s\n", method.Name)
		fmt.Println("Inputs:")
		for j := 0; j < method.Type.NumIn(); j++ {
			fmt.Printf("\t%s\n", method.Type.In(j))
		}
		fmt.Println("Outputs:")
		for k := 0; k < method.Type.NumOut(); k++ {
			fmt.Printf("\t%s\n", method.Type.Out(k))
		}
		fmt.Println()
	}
}

type INode interface {
	Version() (string, error)
	Txes(height int) ([]Tx, error)
}

type Tx struct {
}

type Node struct {
}

func (n *Node) Version() (string, error) {
	return "1.0.0", nil
}

func (n *Node) Txes(height int) ([]Tx, error) {
	return nil, nil
}
