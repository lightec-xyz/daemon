package circuits

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/consensys/gnark/frontend"
	"github.com/lightec-xyz/daemon/common"
	beacon_header "github.com/lightec-xyz/provers/circuits/beacon-header"
	beacon_header_finality "github.com/lightec-xyz/provers/circuits/beacon-header-finality"
	"github.com/lightec-xyz/provers/circuits/fabric/receipt-proof"
	"github.com/lightec-xyz/provers/circuits/fabric/tx-proof"
	"github.com/lightec-xyz/provers/circuits/redeem"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	proverCommon "github.com/lightec-xyz/provers/common"
	reLightCommon "github.com/lightec-xyz/reLight/circuits/common"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	native_plonk "github.com/consensys/gnark/backend/plonk"
	plonk_bn254 "github.com/consensys/gnark/backend/plonk/bn254"
	"github.com/consensys/gnark/backend/witness"
	"github.com/lightec-xyz/btc_provers/circuits/grandrollup"
	"github.com/lightec-xyz/daemon/logger"
	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
	txineth2 "github.com/lightec-xyz/provers/circuits/tx-in-eth2"
	"github.com/lightec-xyz/reLight/circuits/genesis"
	"github.com/lightec-xyz/reLight/circuits/recursive"
	"github.com/lightec-xyz/reLight/circuits/unit"
	"github.com/lightec-xyz/reLight/circuits/utils"
)

var _ ICircuit = (*Circuit)(nil)

type Circuit struct {
	unit      *unit.Unit
	recursive *recursive.Recursive
	genesis   *genesis.Genesis
	Cfg       *CircuitConfig
	debug     bool
}

func (c *Circuit) Load() error {
	// todo
	return nil
}

func NewCircuit(cfg *CircuitConfig) (*Circuit, error) {
	unitConfig := unit.NewUnitConfig(cfg.DataDir, cfg.SrsDir, cfg.SubDir)
	genesisConfig := genesis.NewGenesisConfig(cfg.DataDir, cfg.SrsDir, cfg.SubDir)
	recursiveConfig := recursive.NewRecursiveConfig(cfg.DataDir, cfg.SrsDir, cfg.SubDir)
	var zkDebug bool
	var err error
	zkDebugEnv := os.Getenv(common.ZkDebugEnv)
	if zkDebugEnv != "" {
		zkDebug, err = strconv.ParseBool(zkDebugEnv)
		if err != nil {
			return nil, err
		}
	}
	return &Circuit{
		unit:      unit.NewUnit(unitConfig),
		recursive: recursive.NewRecursive(recursiveConfig),
		genesis:   genesis.NewGenesis(genesisConfig),
		Cfg:       cfg,
		debug:     zkDebug, // todo
	}, nil
}

func (c *Circuit) RedeemProve(txProof, txWitness, bhProof, bhWitness, bhfProof, bhfWitness, beginId, endId, genesisScRoot,
	currentSCSSZRoot []byte, txVar *[tx.MaxTxUint128Len]frontend.Variable, receiptVar *[receipt.MaxReceiptUint128Len]frontend.Variable) (*reLightCommon.Proof, error) {
	//todo
	logger.Debug("current zk circuit RedeemProve")
	if c.debug {
		logger.Warn("current zk circuit RedeemProve prove is debug,skip prove")
		return debugProof()
	}
	txVk, err := utils.ReadVk(filepath.Join(c.Cfg.DataDir, proverCommon.TxInEth2VkFile))
	if err != nil {
		logger.Error("read vk error:%v", err)
		return nil, err
	}

	txInEth2Proof, err := ParseProof(txProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	txInEth2Witness, err := ParseWitness(txWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}

	bhVk, err := utils.ReadVk(filepath.Join(c.Cfg.DataDir, proverCommon.BeaconHeaderOuterVkFile))
	if err != nil {
		logger.Error("read vk error:%v", err)
		return nil, err
	}
	blockHeaderProof, err := ParseProof(bhProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	blockHeaderWitness, err := ParseWitness(bhWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	bhfVk, err := utils.ReadVk(filepath.Join(c.Cfg.DataDir, proverCommon.BeaconHeaderFinalityPkFile))
	if err != nil {
		logger.Error("read vk error:%v", err)
		return nil, err
	}
	blockHeaderFinalityProof, err := ParseProof(bhfProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}

	blockHeaderFinalityWitness, err := ParseWitness(bhfWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	genesisSCSSZRoot, err := GetGenesisSCSSZRoot(genesisScRoot)
	if err != nil {
		logger.Error("get genesis scssz root error:%v", err)
		return nil, err
	}

	proof, err := redeem.Prove(c.Cfg.DataDir, txVk, txInEth2Proof, txInEth2Witness, bhVk, blockHeaderProof, blockHeaderWitness,
		bhfVk, blockHeaderFinalityProof, blockHeaderFinalityWitness, beginId, endId, genesisSCSSZRoot, currentSCSSZRoot, *txVar,
		*receiptVar)
	if err != nil {
		logger.Error("redeem prove error:%v", err)
		return nil, err
	}
	return &reLightCommon.Proof{ // todo
		Proof: proof.Proof,
		Wit:   proof.Wit,
	}, nil
}

func (c *Circuit) BeaconHeaderFinalityUpdateProve(genesisSCSSZRoot string, recursiveProof, recursiveWitness, outerProof,
	outerWitness []byte, finalityUpdate *proverType.FinalityUpdate, scUpdate *proverType.SyncCommitteeUpdate) (*reLightCommon.Proof, error) {
	// todo
	ok, err := common.VerifyLightClientUpdate(scUpdate)
	if err != nil {
		logger.Error("verify light client update error:%v", err)
		return nil, err
	}
	if !ok {
		logger.Error("verify light client update error")
		return nil, fmt.Errorf("verify light client update error")
	}
	logger.Debug("current zk circuit BeaconHeaderFinalityUpdateProve")
	if c.debug {
		logger.Warn("current zk circuit BeaconHeaderFinalityUpdateProve prove is debug,skip prove")
		return debugProof()
	}
	scRecursiveVk, err := utils.ReadVk(filepath.Join(c.Cfg.DataDir, reLightCommon.RecursiveVkFile))
	if err != nil {
		logger.Error("read vk error:%v", err)
		return nil, err
	}
	scRecursiveProof, err := ParseProof(recursiveProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	scRecursiveWitness, err := ParseWitness(recursiveWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	scOuterVk, err := utils.ReadVk(filepath.Join(c.Cfg.DataDir, reLightCommon.OuterVkFile))
	if err != nil {
		logger.Error("read vk error:%v", err)
		return nil, err
	}
	scOuterProof, err := ParseProof(outerProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	scOuterWitness, err := ParseWitness(outerWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}

	proof, err := beacon_header_finality.Prove(c.Cfg.DataDir, genesisSCSSZRoot, scRecursiveVk, scRecursiveProof,
		scRecursiveWitness, scOuterVk, scOuterProof, scOuterWitness, finalityUpdate, scUpdate)
	if err != nil {
		logger.Error("beacon header finality update prove error:%v", err)
		return nil, err
	}
	return &reLightCommon.Proof{ // todo
		Proof: proof.Proof,
		Wit:   proof.Wit,
	}, nil
}

func (c *Circuit) BeaconHeaderProve(header proverType.BeaconHeaderChain) (*reLightCommon.Proof, error) {
	//todo
	logger.Debug("current zk circuit BeaconHeaderProve")
	if c.debug {
		logger.Warn("current zk circuit BeaconHeaderProve prove is debug,skip prove")
		return debugProof()
	}
	proof, err := beacon_header.Prove(c.Cfg.DataDir, header)
	if err != nil {
		logger.Error("beacon header prove error:%v %v %v", header.BeginSlot, header.EndSlot, err)
		return nil, err
	}
	return &reLightCommon.Proof{ // todo
		Proof: proof.Proof,
		Wit:   proof.Wit,
	}, nil
}

func (c *Circuit) TxInEth2Prove(param *ethblock.TxInEth2ProofData) (*reLightCommon.Proof, error) {
	logger.Debug("current zk circuit TxInEth2Prove")
	if c.debug {
		logger.Warn("current zk circuit TxInEth2Prove prove is debug,skip prove")
		return debugProof()
	}
	proof, err := txineth2.Prove(c.Cfg.DataDir, param)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return &reLightCommon.Proof{ // todo
		Proof: proof.Proof,
		Wit:   proof.Wit,
	}, err
}

func (c *Circuit) DepositProve(txId, blockHash string) (*reLightCommon.Proof, error) {
	logger.Debug("current zk circuit DepositProve")
	if c.debug {
		logger.Warn("current zk circuit DepositProve is debug,skip prove ")
		return debugProof()
	}
	return grandrollup.ProveWithDefaults(c.Cfg.DataDir, txId, blockHash)
}

func (c *Circuit) UnitProve(period uint64, update *utils.LightClientUpdateInfo) (*reLightCommon.Proof, *reLightCommon.Proof, error) {
	// todo
	logger.Debug("current zk circuit unit prove")
	ok, err := common.VerifyLightClientUpdate(update)
	if err != nil {
		logger.Error("verify light client update error:%v", err)
		return nil, nil, err
	}
	if !ok {
		logger.Error("verify light client update error")
		return nil, nil, fmt.Errorf("verify light client update error")
	}
	if c.debug {
		logger.Warn("current zk circuit unit prove is debug mode,skip prove")
		proof, err := debugProof()
		if err != nil {
			logger.Error("debug proof error:%v", err)
			return nil, nil, err
		}
		return proof, proof, err
	}
	logger.Warn("really do unit prove now: %v", period)
	subDir := fmt.Sprintf("sc%d", period)
	err = innerProve(c.Cfg.DataDir, subDir, update)
	if err != nil {
		logger.Error("inner prove error:%v", err)
		return nil, nil, err
	}
	outerProof, err := outerProve(c.Cfg.DataDir, subDir, update)
	if err != nil {
		logger.Error("outer prove error:%v", err)
		return nil, nil, err
	}
	unitProof, err := unitProv(c.Cfg.DataDir, subDir, update)
	if err != nil {
		logger.Error("unit prove error:%v", err)
		return nil, nil, err
	}
	return unitProof, outerProof, nil
}

func (c *Circuit) RecursiveProve(choice string, firstProof, secondProof, firstWitness, secondWitness []byte,
	beginId, relayId, endId []byte) (*reLightCommon.Proof, error) {
	logger.Debug("recursive prove request data choice:%v", choice)
	if c.debug {
		logger.Warn("current zk circuit recursive prove is debug mode,skip prove")
		return debugProof()
	}
	if !(choice == "genesis" || choice == "recursive") {
		return nil, fmt.Errorf("invalid choice: %s", choice)
	}
	firstPr, err := ParseProof(firstProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	secondPr, err := ParseProof(secondProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	firstWit, err := ParseWitness(firstWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	secondWit, err := ParseWitness(secondWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	config := recursive.NewRecursiveConfig(c.Cfg.DataDir, c.Cfg.SrsDir, "")
	recursiveCir := recursive.NewRecursive(config)
	err = recursiveCir.Load()
	if err != nil {
		logger.Error("recursive load error:%v", err)
		return nil, err
	}
	proof, err := recursiveCir.Prove(choice, firstPr, secondPr, firstWit, secondWit, beginId, relayId, endId)
	if err != nil {
		logger.Error("recursive prove error:%v", err)
		return nil, err
	}
	return proof, err
}

func (c *Circuit) GenesisProve(firstProof, secondProof, firstWitness, secondWitness []byte,
	genesisId, firstId, secondId []byte) (*reLightCommon.Proof, error) {
	logger.Debug("current zk circuit syncCommittee genesis prove")
	if c.debug {
		logger.Warn("current zk circuit genesis prove is debug mode,skip prove")
		return debugProof()
	}
	firstPf, err := ParseProof(firstProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	secondPf, err := ParseProof(secondProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	firstWit, err := ParseWitness(firstWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	secondWit, err := ParseWitness(secondWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	config := genesis.NewGenesisConfig(c.Cfg.DataDir, "", "")
	genesisCir := genesis.NewGenesis(config)
	err = genesisCir.Load()
	if err != nil {
		logger.Error("genesis load error:%v", err)
		return nil, err
	}
	proof, err := genesisCir.Prove(firstPf, secondPf, firstWit, secondWit, genesisId, firstId, secondId)
	if err != nil {
		logger.Error("genesis prove error:%v", err)
		return nil, err
	}
	return proof, err
}

func (c *Circuit) UpdateChangeProve(txId, blockHash string) (*reLightCommon.Proof, error) {
	logger.Debug("current zk circuit UpdateChangeProve")
	if c.debug {
		logger.Warn("current zk circuit DepositProve is debug,skip prove ")
		return debugProof()
	}
	return grandrollup.ProveWithDefaults(c.Cfg.DataDir, txId, blockHash)
}

func SyncCommitRoot(update *utils.LightClientUpdateInfo) ([]byte, error) {
	ok, err := common.VerifyLightClientUpdate(update)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("verify light client update error")
	}
	return utils.SyncCommitRoot(update)
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

func unitProv(dataDir string, subDir string, update *utils.LightClientUpdateInfo) (*reLightCommon.Proof, error) {
	unitCfg := unit.NewUnitConfig(dataDir, "", subDir)
	unit := unit.NewUnit(unitCfg)
	err := unit.Load()
	if err != nil {
		return nil, err
	}
	proofs, err := unit.Prove(update)
	if err != nil {
		return nil, err
	}
	return proofs, nil
}

func outerProve(dataDir string, subDir string, update *utils.LightClientUpdateInfo) (*reLightCommon.Proof, error) {
	outerCfg := unit.NewOuterConfig(dataDir, "", subDir)
	outer := unit.NewOuter(&outerCfg)
	err := outer.Load()
	if err != nil {
		return nil, err
	}
	proofs, err := outer.Prove(update)
	if err != nil {
		return nil, err
	}
	err = outer.Save(proofs)
	if err != nil {
		return nil, err
	}
	return proofs, nil
}

func innerProve(dataDir string, subDir string, update *utils.LightClientUpdateInfo) error {
	innerCfg := unit.NewInnerConfig(dataDir, "", subDir)
	inner := unit.NewInner(&innerCfg)
	err := inner.Load()
	if err != nil {
		return err
	}
	assignments, err := inner.GetCircuitAssignments(update)
	if err != nil {
		return err
	}
	for index, assignment := range assignments {
		proof, err := inner.Prove(assignment)
		if err != nil {
			return err
		}
		err = inner.Save(index, proof)
		if err != nil {
			return err
		}
	}
	return nil
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
	DataDir string
	SrsDir  string
	SubDir  string
	Debug   bool
}

func debugProof() (*reLightCommon.Proof, error) {
	// todo only just local debug
	time.Sleep(15 * time.Second)
	field := ecc.BN254.ScalarField()
	w, err := witness.New(field)
	if err != nil {
		return nil, err
	}
	witnessByts, err := hex.DecodeString("000000180000000000000018000000000000000000000000000000000000000000000000bc4d9a773a304f7c000000000000000000000000000000000000000000000000c879892de7b1130b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004d13e6221265d5470000000000000000000000000000000000000000000000000a9a955cdf54319900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000bd806aa2440faf3a00000000000000000000000000000000000000000000000056bb0ec865d27e9800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000fba8545ab164e9ef0000000000000000000000000000000000000000000000000653e66962364b88000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007e1da36c41365c0d000000000000000000000000000000000000000000000000942a5884da9b98da00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000061ae6bd87e134e80000000000000000000000000000000000000000000000002085a29e1cf057bb00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		return nil, err
	}
	_, err = w.ReadFrom(bytes.NewReader(witnessByts))
	if err != nil {
		return nil, err
	}
	proofBytes, err := hex.DecodeString("d05f24ef1d3e8b59ba060b4141b1048b0199ed002c0a313c03b87f8e51e24aabd01a1fb0eaee9473970eeecaf7c3c850b57afe3386641bcc6cc18581b003abeacf557f98fc31079075dccf75f355221a018edd7bf5054fdbd280f3c61f537be6e3cba22a1f3ba940149f4f251195cea82b6559344ae7d6d40e600655cf591561cc8422fbc1ba2929de9488315e5d23d717e8c9048d534d3569a358f57eddfcb5dfa8e221f4104d063048df28114ac4ab5a7883245b55901367b972b4ead270c0a6550b7c6c44fd672a5b88ed8d153cf50eb6247d2bf48794ecb3803a2017e967d3553de5efb7ca588f31ed43ce43f198619d6eecc1203970caab1f46123b8b520000000805994c4caf545b0998b1cd70a2778274a35f5b8a2d1c64344b7b119b62f990232164704f0f9cd6e2ea6de50ca2694790cbd6c5db0a14ac6d4462f0563f5d42e120e2fa2be493efd68b25a793957779ab2af40f2b0422e18ba72bdd57196c81b026b9f7955a08ee21bea6045e64eeca6cf7e7504b290960e5ecd58b1919ced66625c08a8c391b0763abbcd3f5ef0509d445ec02f3b11db660796ae1d02f47f0b91bdf22f73b1b9fb3e7a8264523bc016164aea7770b7f47223a4c449833f9324b217ce713ec851098916ce7b9349ee7bfc63095e19644ce496e8cdd542d0703aa0723b8b49eae51db612335d883c26d2013663ef4c3fb8ed734afeef30bd38e19a5415f3c287648a2160c32b1176ffd043d27fa50614843649306d814b9f198a70cfa4c03f0f487ad2a8cd3f0d3be71cccc00d237eba86e4f9d5ed91cfc0da7f500000001840fdac67c39e3ccf5363c17dca27f6118bdc2dd629e07885be3778401fe566c")
	if err != nil {
		return nil, err
	}
	proof := &plonk_bn254.Proof{}
	_, err = proof.ReadFrom(bytes.NewReader(proofBytes))
	if err != nil {
		return nil, err
	}
	return &reLightCommon.Proof{
		Proof: proof,
		Wit:   w,
	}, nil
}

func GetGenesisSCSSZRoot(root []byte) ([2]frontend.Variable, error) {
	panic(root)
}

func GetTxVar(data []byte) ([tx.MaxTxUint128Len]frontend.Variable, error) {
	panic(data)
}

func GetReceiptVar(data []byte) ([receipt.MaxReceiptUint128Len]frontend.Variable, error) {
	panic(data)
}

func ProofToHexSol(proof native_plonk.Proof) (string, error) {
	_proof := proof.(*plonk_bn254.Proof)
	proofStr := hex.EncodeToString(_proof.MarshalSolidity())
	return proofStr, nil

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
