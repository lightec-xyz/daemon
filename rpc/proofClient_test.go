package rpc

import "testing"

func TestProofClient_GenZkProof(t *testing.T) {
	client, err := NewWsProofClient("ws://127.0.0.1:30001")
	if err != nil {
		t.Fatal(err)
	}
	maxNums, err := client.MaxNums()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(maxNums)
	currentNums, err := client.CurrentNums()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(currentNums)
	proofInfo, err := client.ProofInfo("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(proofInfo)
	depositProof, err := client.TxInEth2Prove(&TxInEth2ProveRequest{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(depositProof)
	redeemProof, err := client.GenRedeemProof(&RedeemRequest{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(redeemProof)

	verifyProof, err := client.TxInEth2Prove(&TxInEth2ProveRequest{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(verifyProof)

	genesisProof, err := client.GenSyncCommGenesisProof(SyncCommGenesisRequest{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(genesisProof)
	unitProof, err := client.GenSyncCommitUnitProof(SyncCommUnitsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(unitProof)
	recursiveProof, err := client.GenSyncCommRecursiveProof(SyncCommRecursiveRequest{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(recursiveProof)

}
