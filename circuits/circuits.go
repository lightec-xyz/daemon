package circuits

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	native_plonk "github.com/consensys/gnark/backend/plonk"
	plonk_bn254 "github.com/consensys/gnark/backend/plonk/bn254"
	"github.com/consensys/gnark/backend/witness"
	"github.com/lightec-xyz/btc_provers/circuits/grandrollup"
	dCom "github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
	txineth2 "github.com/lightec-xyz/provers/circuits/tx-in-eth2"
	"github.com/lightec-xyz/reLight/circuits/common"
	"github.com/lightec-xyz/reLight/circuits/genesis"
	"github.com/lightec-xyz/reLight/circuits/recursive"
	"github.com/lightec-xyz/reLight/circuits/unit"
	"github.com/lightec-xyz/reLight/circuits/utils"
)

type Circuit struct {
	unit      *unit.Unit
	recursive *recursive.Recursive
	genesis   *genesis.Genesis
	Cfg       *CircuitConfig
	debug     bool
}

func NewCircuit(cfg *CircuitConfig) (*Circuit, error) {
	unitConfig := unit.NewUnitConfig(cfg.DataDir, cfg.SrsDir, cfg.SubDir)
	genesisConfig := genesis.NewGenesisConfig(cfg.DataDir, cfg.SrsDir, cfg.SubDir)
	recursiveConfig := recursive.NewRecursiveConfig(cfg.DataDir, cfg.SrsDir, cfg.SubDir)
	var zkDebug bool
	var err error
	zkDebugEnv := os.Getenv(dCom.ZkDebugEnv)
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

func (c *Circuit) Load() error {
	// todo
	return nil
	//if c.debug {
	//	logger.Warn("current zk circuit is debug mode,skip load")
	//	return nil
	//}
	//// todo
	//err := c.genesis.Load()
	//if err != nil {
	//	logger.Error("genesis load error:%v", err)
	//	return err
	//}
	//err = c.unit.Load()
	//if err != nil {
	//	logger.Error("unit load error:%v", err)
	//	return err
	//}
	//err = c.recursive.Load()
	//if err != nil {
	//	logger.Error("recursive load error:%v", err)
	//	return err
	//}
	//return nil
}

func (c *Circuit) TxInEth2Prove(param *ethblock.TxInEth2ProofData) (*common.Proof, error) {
	if c.debug {
		logger.Warn("current zk circuit TxInEth2Prove prove is debug,skip prove")
		return debugProof()
	}
	proof, err := txineth2.Prove(c.Cfg.DataDir, param)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return proof, err
}

func (c *Circuit) DepositProve(txId, blockHash string) (*common.Proof, error) {
	if c.debug {
		logger.Warn("current zk circuit DepositProve is debug,skip prove ")
		return debugProof()
	}
	return grandrollup.ProveWithDefaults(c.Cfg.DataDir, txId, blockHash)
}

func (c *Circuit) UnitProve(period uint64, update *utils.LightClientUpdateInfo) (*common.Proof, error) {
	if c.debug {
		logger.Warn("current zk circuit unit prove is debug mode,skip prove")
		return debugProof()
	}
	// todo
	logger.Warn("really do unit prove now: %v", period)
	//proof, err := unitProve(c.Cfg.DataDir, c.Cfg.SrsDir, fmt.Sprintf("sc%d", period), update)
	//proof, err := c.unit.Prove(update)
	subDir := fmt.Sprintf("sc%d", period)
	err := innerProve(c.Cfg.DataDir, subDir, update)
	if err != nil {
		logger.Error("inner prove error:%v", err)
		return nil, err
	}
	err = outerProve(c.Cfg.DataDir, subDir, update)
	if err != nil {
		logger.Error("outer prove error:%v", err)
		return nil, err
	}
	proof, err := innerUnitProv(c.Cfg.DataDir, subDir, update)
	if err != nil {
		logger.Error("unit prove error:%v", err)
		return nil, err
	}
	return proof, nil
}

func (c *Circuit) RecursiveProve(choice string, firstProof, secondProof, firstWitness, secondWitness []byte,
	beginId, relayId, endId []byte) (*common.Proof, error) {
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
	genesisId, firstId, secondId []byte) (*common.Proof, error) {
	//logger.Debug("genesis prove request data firstProof:%x secondProof:%x firstWitness:%x secondWitness:%x,genesisId:%x firstId:%x secondId:%x",
	//	firstProof, secondProof, firstWitness, secondWitness, genesisId, firstId, secondId)
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
func (c *Circuit) TxBlockIsParentOfCheckPointProve() (*common.Proof, error) {

	return nil, nil
}

func (c *Circuit) FinalityUpdateProve() (*common.Proof, error) {

	return nil, nil
}

func (c *Circuit) RedeemProve() (*common.Proof, error) {
	panic(c)
	return nil, nil
}
func (c *Circuit) UpdateChangeProve(txId, blockHash string) (*common.Proof, error) {
	if c.debug {
		logger.Warn("current zk circuit DepositProve is debug,skip prove ")
		return debugProof()
	}
	return grandrollup.ProveWithDefaults(c.Cfg.DataDir, txId, blockHash)
}

func SyncCommitRoot(update *utils.LightClientUpdateInfo) ([]byte, error) {
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

func innerUnitProv(dataDir string, subDir string, update *utils.LightClientUpdateInfo) (*common.Proof, error) {
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
	if err != nil {
		return nil, err
	}
	return proofs, nil
}

func outerProve(dataDir string, subDir string, update *utils.LightClientUpdateInfo) error {
	outerCfg := unit.NewOuterConfig(dataDir, "", subDir)
	outer := unit.NewOuter(&outerCfg)
	err := outer.Load()
	if err != nil {
		return err
	}
	proofs, err := outer.Prove(update)
	if err != nil {
		return err
	}
	err = outer.Save(proofs)
	if err != nil {
		return err
	}
	return nil
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

func debugProof() (*common.Proof, error) {
	// todo only just local debug
	time.Sleep(15 * time.Second)
	field := ecc.BN254.ScalarField()
	w, err := witness.New(field)
	if err != nil {
		return nil, err
	}
	//proofBytes, err := hex.DecodeString("201c5d2f7f746c025582d4bf687c821d46ef2cf7d79b18a648731132aad3b1402b171c26a4b247cb749babdd443cf916774a4b74e85e3ea9ba19f4ee43638dd7032c6ea702809b843263e19d4b97cfdf2471a74a127b56ad581473c6df3d0f472d2eea103495c5e856bf5aee8b3270bd9adab84209338b6f09f51cf39a365c351ada496af014c9abf7276d854ea0d6ea6a36e090d0a45b9069f81fcb6f9fbd73032bb31aa39995a864d98c440a77b8950847066150f0e94f066054bae18704811680a1a766672149d412d5b60c57b31433f897942bf7216aa3aaf1642e4992e60c007447583c5ce87e207476241fcca408417aeff488f7165f60bce73ce7e12a171f8927ca72e68714e0a80c89f833fdbdf895cbe94e6e6892cd474d1964970d2313ad60cffff9450d1f9724fc0bb68d81f031a09209406232d5f4d7dfb4c15e262822a536f433d56b461c812916ecf692a979a0246fc8de7de0d3092c7e7a8e283342712db0f2db4ad60551e8e110799ea7bafd37990658473c9e447f93cd8809bb45473b7d249831878bd88381c76e3695ab957ff876f5fb75c3091ae7a2c82877b98e4aaca1ae32064360f9ad05d0f1595d5613a416ec544192309b8ed61b21ef0de29b4706bb0f0e91a82945168ebbe619b6987e92790a9920444de540e0264525ca3ac9f4de4a8409226e52bb18ce38d837c700f7aee5a907b45a0b5fcb06c9f5d06cb0a8558e239b09f28c90de463dc2c9d12cb985d938cca3e4cd459825bbe0d37f3de8e6a212a181cb3f5b25320afbe17a1aab305bd55e0fdaaf4a851d0cc3ae31c567fb9b466e9049b9db84936191b9cfa7f709e269c40982cfb5ca001915a0707f2a0484a92cc9befdf16e5213972d39e1cb7b2cf421159f2836c21941febd827e29c1b669181c00a9cee546cb1ee5eb391d6841f4dff4d5acc455058d495e2a39d5fc9a1a1902c7b539c791f7c72301391bac7a127a5dadbb289c16851ed31eb91037cb6f61b8d9628af08983cf94273b7425a0a96995103c95ee267b2a4d36a83faafde99eaa295de72cf58defd0e9482d4979dec70aa10e4fe506a0ab76dbb2ae094f758e74aa34d7c433fcc70bc88c3f12a878881ff6c0fee81f759700345b16aea15f43ee6ab7b3803c4d149b667ae8570d21783e1cc294dc06b40f01e8d5a42ecd7f6d50bf0b3209054a2acf7d27791754445a87a1baace024f3b5db24c273a368285f93531bba3693df99680cb56079117d301988dc2731087738b148b8626517c5768a6ad37eeb6a29c3f83145b14572cf46e29aae0660")
	//if err != nil {
	//	return nil, err
	//}
	//buffer := bytes.NewBuffer(proofBytes)
	proof := &plonk_bn254.Proof{}
	//_, err = proof.WriteTo(buffer)
	//if err != nil {
	//	return nil, err
	//}
	return &common.Proof{
		Proof: proof,
		Wit:   w,
	}, nil
}

func ProofToHexSolBytes(proof native_plonk.Proof) (string, error) {
	_proof := proof.(*plonk_bn254.Proof)
	proofStr := hex.EncodeToString(_proof.MarshalSolidity())
	return proofStr, nil

}

func ProofToBytes(proof native_plonk.Proof) []byte {
	var buf bytes.Buffer
	_, err := proof.WriteTo(&buf)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

func WitnessToBytes(witness witness.Witness) []byte {
	var buf bytes.Buffer
	_, err := witness.WriteTo(&buf)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}
