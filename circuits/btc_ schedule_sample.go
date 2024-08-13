package circuits

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/lightec-xyz/btc_provers/circuits/common"
)

var genesisBlockHeight, cpBlockHeight uint32
var ic ICircuit

func schedule(latestBlockHeight uint32) {
	var blockhash chainhash.Hash
	var depositTxs, redeemTxs []chainhash.Hash
	var txBlockHeight uint32

	blockhash, txBlockHeight, depositTxs, redeemTxs := ScanBlock(latestBlockHeight - 1)

	// for tx
	countDeposit, countRedeem := len(depositTxs), len(redeemTxs)
	if countRedeem != 0 || countDeposit != 0 {
		// for chain
		ic.BtcBaseProve()
		ic.BtcChainProve()

		// for cp depth, usually >= MinPacked
		ic.BtcBulkProve(latestBlockHeight, cpBlockHeight)
		if latestBlockHeight-cpBlockHeight >= common.MinPacked {
			ic.BtcPackProve(latestBlockHeight, cpBlockHeight)
		}

		// for tx depth, usually < MinPacked
		ic.BtcBulkProve(latestBlockHeight, txBlockHeight)
		if latestBlockHeight-txBlockHeight >= common.MinPacked {
			ic.BtcPackProve(latestBlockHeight, txBlockHeight)
		}

		if countRedeem != 0 {
			ic.BtcChangeProve()
		}

		if countDeposit != 0 {
			ic.BtcDepositProve()
		}
	}

	// for chain recursive
	if (latestBlockHeight-genesisBlockHeight+1)%common.CapacityBaseLevel == 0 {
		ic.BtcBaseProve()
		ic.BtcMiddleProve()
	}
	if (latestBlockHeight-genesisBlockHeight+1)%common.CapacitySuperBatch == 0 {
		ic.BtcUpperProve()
	}
	if (latestBlockHeight-genesisBlockHeight+1)%common.CapacityDifficultyBlock == 0 {
		ic.BtcDuperRecursiveProve()
	}

	// for cp depth recursive
	if (latestBlockHeight-cpBlockHeight)%common.CapacityBulkUint == 0 {
		ic.BtcBulkProve()
		ic.BtcDepthRecursiveProve()
	}
}
