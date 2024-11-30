package resulthandlers

import (
	"intellix/x/price/dvs/types"
	"math/big"

	contractDataOracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"golang.org/x/crypto/sha3"
)

type ProcessRequestPriceFeedResultHandler struct {
}

func NewProcessRequestPriceFeedResultHandler() *ProcessRequestPriceFeedResultHandler {
	return &ProcessRequestPriceFeedResultHandler{}
}

// TODO: put it in a common location
func packUint256(value *big.Int) ([]byte, error) {
	uint256Type, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return nil, err
	}

	arguments := abi.Arguments{{Type: uint256Type}}

	return arguments.Pack(value)
}

func (p *ProcessRequestPriceFeedResultHandler) getAbiEncodeData(msg proto.Message) ([]byte, error) {
	r, ok := msg.(*types.RequestPriceFeedOut)
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
			Name: "data",
			Type: "bytes",
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

	packedPrice, err := packUint256(r.Price.BigInt())
	if err != nil {
		return nil, err
	}

	bytes, err := arguments.Pack(&contractDataOracle.IDataOracleTaskResponse{
		ReferenceTaskIndex: r.TaskIndex,
		Data:               packedPrice,
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
