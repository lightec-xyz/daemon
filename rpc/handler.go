package rpc

import (
	"fmt"
)

type Handler struct {
}

func (h *Handler) HelloWorld(name string) (string, error) {
	fmt.Printf("req: %v", name)
	return fmt.Sprintf("%v resp,age: %v", name, "dd"), nil
}
