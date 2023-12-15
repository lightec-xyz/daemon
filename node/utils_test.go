package node

import (
	"fmt"
	"testing"
)

func TestFloatToBig(t *testing.T) {
	result := fmt.Sprintf("%0.8f", 0.123456789)
	t.Log(result)
}
