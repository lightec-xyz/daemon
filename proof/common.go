package proof

import (
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
)

func toZkProofType(proofTypes []string) ([]common.ZkProofType, error) {
	var zkProofTypes []common.ZkProofType
	for _, proofType := range proofTypes {
		ptype, err := common.ToZkProofType(proofType)
		if err != nil {
			logger.Error("convert proof type error:%v %v", proofType, err)
			return nil, err
		}
		zkProofTypes = append(zkProofTypes, ptype)
	}
	return zkProofTypes, nil
}
