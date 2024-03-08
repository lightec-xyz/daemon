package oasis

import (
	"context"
	"github.com/oasisprotocol/oasis-core/go/common"
	"github.com/oasisprotocol/oasis-sdk/client-sdk/go/client"
	"google.golang.org/grpc"
	"testing"
)

func TestDemo(t *testing.T) {
	clientConn, err := grpc.Dial("127.0.0.1:9090", grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	var runtimeID common.Namespace
	runtimeClient := client.New(clientConn, runtimeID)
	err = client.NewTransactionBuilder(runtimeClient, "", nil).SubmitTx(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

}
