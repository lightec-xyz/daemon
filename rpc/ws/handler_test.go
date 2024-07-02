package ws

import "testing"

func TestHandler(t *testing.T) {
	node := Node{}
	handler := newHandler(&node)
	result := handler.call("Version")
	t.Log(result)
	t.Log(node.Count)
	result = handler.call("Version")
	t.Log(node.Count)

}
