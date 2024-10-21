package resulthandlers

import (
	contractPriceOracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/PriceOracle"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"golang.org/x/crypto/sha3"
	"intellix/x/price/dvs/types"
)

type ProcessRequestPriceFeedResultHandler struct {
}

func NewProcessRequestPriceFeedResultHandler() *ProcessRequestPriceFeedResultHandler {
	return &ProcessRequestPriceFeedResultHandler{}
}

func (p *ProcessRequestPriceFeedResultHandler) getAbiEncodeData(msg proto.Message) ([]byte, error) {
	r, ok := msg.(*types.AggregatedRequestPrice)
	if !ok {
		return nil, nil
	}

	// The order here has to match the field ordering of cstaskmanager.IIncrediblePriceTaskManagerV2TaskResponse
	taskResponseType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "referenceTaskIndex",
			Type: "uint32",
		},
		{
			Name: "price",
			Type: "uint256",
		},
	})
	if err != nil {
		return nil, err
	}
	arguments := abi.Arguments{
		{
			Type: taskResponseType,
		},
	}

	bytes, err := arguments.Pack(&contractPriceOracle.IPriceOracleTaskResponse{
		ReferenceTaskIndex: r.TaskIndex,
		Price:              r.Price.BigInt(),
	})
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (p *ProcessRequestPriceFeedResultHandler) GetData(msg proto.Message) ([]byte, error) {
	return p.getAbiEncodeData(msg)
}

func (p *ProcessRequestPriceFeedResultHandler) GetDigest(msg proto.Message) ([]byte, error) {
	data, err := p.getAbiEncodeData(msg)
	if err != nil {
		return nil, err
	}
	// calc digest
	var taskResponseDigest [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(data)
	copy(taskResponseDigest[:], hasher.Sum(nil)[:32])

	return taskResponseDigest[:], nil
}
