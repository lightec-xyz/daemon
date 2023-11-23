package bitcoin

import (
	"fmt"
	"testing"
)

func TestClient_scantxoutset(t *testing.T) {

	result, err := client.Scantxoutset("bcrt1q823m29mq29lt7heewwdry34phtgu767lgv5mx6")
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

}
