package rpc

import "testing"

func TestProofClient_GenZkProof(t *testing.T) {
	client, err := NewWsProofClient("ws://127.0.0.1:30001")
	if err != nil {
		t.Fatal(err)
	}
	proofInfo, err := client.ProofInfo("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(proofInfo)
	depositProof, err := client.GenDepositProof(DepositRequest{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(depositProof)
	redeemProof, err := client.GenRedeemProof(RedeemRequest{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(redeemProof)

	verifyProof, err := client.GenVerifyProof(VerifyRequest{})
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
