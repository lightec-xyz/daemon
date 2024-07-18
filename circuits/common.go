package circuits

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	native_plonk "github.com/consensys/gnark/backend/plonk"
	plonk_bn254 "github.com/consensys/gnark/backend/plonk/bn254"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	btcprovercom "github.com/lightec-xyz/btc_provers/circuits/common"
	"github.com/lightec-xyz/chainark"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	reLightCommon "github.com/lightec-xyz/reLight/circuits/common"
	"github.com/lightec-xyz/reLight/circuits/utils"
	"math/big"
	"time"
)

func SyncCommitRoot(update *utils.SyncCommitteeUpdate) ([]byte, error) {
	ok, err := common.VerifyLightClientUpdate(update)
	if err != nil {
		logger.Error("verify light client update error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("verify light client update error")
	}
	return utils.SyncCommitRoot(update)
}

func ParseWitness(body []byte) (witness.Witness, error) {
	field := ecc.BN254.ScalarField()
	reader := bytes.NewReader(body)
	wit, err := witness.New(field)
	if err != nil {
		return nil, err
	}
	_, err = wit.ReadFrom(reader)
	if err != nil {
		return nil, err
	}
	return wit, nil
}

type CircuitConfig struct {
	DataDir  string
	SrsDir   string
	SetupDir string
	Debug    bool
}

func debugProof() (*reLightCommon.Proof, error) {
	// todo just test for debug
	time.Sleep(5 * time.Second)
	field := ecc.BN254.ScalarField()
	w, err := witness.New(field)
	if err != nil {
		logger.Error("new witness error:%v", err)
		return nil, err
	}
	witnessByts, err := hex.DecodeString("000000180000000000000018000000000000000000000000000000000000000000000000bc4d9a773a304f7c000000000000000000000000000000000000000000000000c879892de7b1130b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004d13e6221265d5470000000000000000000000000000000000000000000000000a9a955cdf54319900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000bd806aa2440faf3a00000000000000000000000000000000000000000000000056bb0ec865d27e9800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000fba8545ab164e9ef0000000000000000000000000000000000000000000000000653e66962364b88000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007e1da36c41365c0d000000000000000000000000000000000000000000000000942a5884da9b98da00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000061ae6bd87e134e80000000000000000000000000000000000000000000000002085a29e1cf057bb00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		logger.Error("hex decode error:%v", err)
		return nil, err
	}
	_, err = w.ReadFrom(bytes.NewReader(witnessByts))
	if err != nil {
		logger.Error("read witness error:%v", err)
		return nil, err
	}
	proofBytes, err := hex.DecodeString("d05f24ef1d3e8b59ba060b4141b1048b0199ed002c0a313c03b87f8e51e24aabd01a1fb0eaee9473970eeecaf7c3c850b57afe3386641bcc6cc18581b003abeacf557f98fc31079075dccf75f355221a018edd7bf5054fdbd280f3c61f537be6e3cba22a1f3ba940149f4f251195cea82b6559344ae7d6d40e600655cf591561cc8422fbc1ba2929de9488315e5d23d717e8c9048d534d3569a358f57eddfcb5dfa8e221f4104d063048df28114ac4ab5a7883245b55901367b972b4ead270c0a6550b7c6c44fd672a5b88ed8d153cf50eb6247d2bf48794ecb3803a2017e967d3553de5efb7ca588f31ed43ce43f198619d6eecc1203970caab1f46123b8b520000000805994c4caf545b0998b1cd70a2778274a35f5b8a2d1c64344b7b119b62f990232164704f0f9cd6e2ea6de50ca2694790cbd6c5db0a14ac6d4462f0563f5d42e120e2fa2be493efd68b25a793957779ab2af40f2b0422e18ba72bdd57196c81b026b9f7955a08ee21bea6045e64eeca6cf7e7504b290960e5ecd58b1919ced66625c08a8c391b0763abbcd3f5ef0509d445ec02f3b11db660796ae1d02f47f0b91bdf22f73b1b9fb3e7a8264523bc016164aea7770b7f47223a4c449833f9324b217ce713ec851098916ce7b9349ee7bfc63095e19644ce496e8cdd542d0703aa0723b8b49eae51db612335d883c26d2013663ef4c3fb8ed734afeef30bd38e19a5415f3c287648a2160c32b1176ffd043d27fa50614843649306d814b9f198a70cfa4c03f0f487ad2a8cd3f0d3be71cccc00d237eba86e4f9d5ed91cfc0da7f500000001840fdac67c39e3ccf5363c17dca27f6118bdc2dd629e07885be3778401fe566c")
	if err != nil {
		logger.Error("hex decode error:%v", err)
		return nil, err
	}
	proof := &plonk_bn254.Proof{}
	_, err = proof.ReadFrom(bytes.NewReader(proofBytes))
	if err != nil {
		logger.Error("read proof error:%v", err)
		return nil, err
	}
	return &reLightCommon.Proof{
		Proof: proof,
		Wit:   w,
	}, nil
}

func HashToLinkageId(hash string) (chainark.LinkageID, error) {
	beginHashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return chainark.LinkageID{}, err
	}
	linkageID := chainark.LinkageIDFromBytes(beginHashBytes, btcprovercom.IdBitsPerVar)
	return linkageID, nil

}

func GetGenesisSCSSZRoot(root string) ([2]frontend.Variable, error) {
	hexBytes, err := HexToBytes(root)
	if err != nil {
		return [2]frontend.Variable{}, err
	}
	if len(hexBytes) < 17 {
		return [2]frontend.Variable{}, fmt.Errorf("genesis scssz root error")
	}
	var genesisSCSSZRoot [2]frontend.Variable
	genesisSCSSZRoot[0] = hexBytes[:16]
	genesisSCSSZRoot[1] = hexBytes[16:]
	return genesisSCSSZRoot, nil
}

func ProofToSolBytes(proof native_plonk.Proof) ([]byte, error) {
	_proof, ok := proof.(*plonk_bn254.Proof)
	if !ok {
		return nil, fmt.Errorf("proof to bn154 error")
	}
	return _proof.MarshalSolidity(), nil
}

func ProofToBytes(proof native_plonk.Proof) ([]byte, error) {
	var buf bytes.Buffer
	_, err := proof.WriteTo(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func WitnessToBytes(witness witness.Witness) ([]byte, error) {
	var buf bytes.Buffer
	pubWit, err := witness.Public()
	if err != nil {
		return nil, err
	}
	_, err = pubWit.WriteTo(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func HexToWitness(witness string) (witness.Witness, error) {
	witnessBytes, err := HexToBytes(witness)
	if err != nil {
		return nil, err
	}
	return ParseWitness(witnessBytes)
}

func ParseProof(proof []byte) (native_plonk.Proof, error) {
	reader := bytes.NewReader(proof)
	var bn254Proof plonk_bn254.Proof
	_, err := bn254Proof.ReadFrom(reader)
	if err != nil {
		return nil, err
	}
	return &bn254Proof, nil
}

func HexToProof(proof string) (native_plonk.Proof, error) {
	proofBytes, err := HexToBytes(proof)
	if err != nil {
		return nil, err
	}
	return ParseProof(proofBytes)
}

func HexToBytes(data string) ([]byte, error) {
	if data[0:2] == "0x" {
		data = data[2:]
	}
	bytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func HexWitnessToBigInts(witness string) ([]*big.Int, error) {
	wit, err := HexToWitness(witness)
	if err != nil {
		return nil, err
	}
	list, ok := wit.Vector().(fr.Vector)
	if !ok {
		return nil, fmt.Errorf("parse fr vector error")
	}
	var bigList []*big.Int
	for _, item := range list {
		value, ok := big.NewInt(0).SetString(item.String(), 10)
		if !ok {
			return nil, fmt.Errorf("parse big int error")
		}
		bigList = append(bigList, value)
	}
	return bigList, nil
}
