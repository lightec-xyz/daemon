package circuits

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	native_plonk "github.com/consensys/gnark/backend/plonk"
	plonk_bn254 "github.com/consensys/gnark/backend/plonk/bn254"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/emulated/sw_bn254"
	"github.com/consensys/gnark/std/recursion/plonk"
	"github.com/lightec-xyz/reLight/circuits/unit"
	"github.com/lightec-xyz/reLight/circuits/utils"
	"os"
	"path/filepath"
)

type Unit struct {
}

func NewUnit() *Unit {
	return &Unit{}
}

func (c *Unit) Verify(opt *OptUnit, update *utils.LightClientUpdateInfo) (bool, error) {
	if opt == nil {
		return false, fmt.Errorf("nil opt")
	}
	if opt.DataDir == "" || opt.SubDir == "" {
		return false, fmt.Errorf("nil DataDir or SubDir")
	}
	dataDir := opt.DataDir
	subDir := opt.SubDir
	vkFile := filepath.Join(dataDir, unit.UnitVkFile)
	proofFile := filepath.Join(dataDir, subDir, fmt.Sprintf("/%v_unit_proof.proof", subDir))
	fproof, err := os.Open(proofFile)
	if err != nil {
		return false, err
	}
	defer fproof.Close()

	var bn254Proof plonk_bn254.Proof
	_, err = bn254Proof.ReadFrom(fproof)
	if err != nil {
		return false, err
	}

	innerField := ecc.BN254.ScalarField()
	outerField := ecc.BN254.ScalarField()

	fvk, err := os.Open(vkFile)
	if err != nil {
		return false, err
	}
	var bn254Vk plonk_bn254.VerifyingKey
	_, err = bn254Vk.UnsafeReadFrom(fvk)
	if err != nil {
		return false, err
	}
	assignment, err := unit.BuildUnitProofCircuitAssignment[sw_bn254.ScalarField, sw_bn254.G1Affine,
		sw_bn254.G2Affine, sw_bn254.GTEl](dataDir, subDir, update)
	if err != nil {
		return false, err
	}
	wit, err := frontend.NewWitness(assignment, innerField)
	if err != nil {
		return false, err
	}
	pubWit, err := wit.Public()
	if err != nil {
		return false, err
	}

	err = native_plonk.Verify(&bn254Proof, &bn254Vk, pubWit, plonk.GetNativeVerifierOptions(outerField, innerField))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *Unit) GenerateProof(opt *OptUnit) error {
	if opt == nil || opt.SrsDataDir == "" || opt.DataDir == "" || opt.SubDir == "" || opt.ParamFile == "" {
		return fmt.Errorf("opt param is empty")
	}
	srsDir := opt.SrsDataDir
	dataDir := opt.DataDir
	err := unit.SetupInnerCircuit(srsDir, dataDir)
	if err != nil {
		return err
	}
	err = unit.SetupOuterCircuit(srsDir, dataDir)
	if err != nil {
		return err
	}
	err = unit.SetupUnitCircuit(srsDir, dataDir)
	if err != nil {
		return err
	}
	err = unit.SetupUnitProofCircuit(srsDir, dataDir)
	if err != nil {
		return err
	}
	err = unit.Prove[sw_bn254.ScalarField, sw_bn254.G1Affine, sw_bn254.G2Affine, sw_bn254.GTEl](dataDir, opt.SubDir, opt.ParamFile)
	if err != nil {
		return err
	}
	return nil
}

type OptUnit struct {
	SrsDataDir string
	DataDir    string
	SubDir     string
	ParamFile  string
}
