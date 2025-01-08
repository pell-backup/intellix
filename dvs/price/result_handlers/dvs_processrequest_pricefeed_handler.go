package resulthandlers

import (
	"intellix/dvs/price/types"
	"intellix/pkg/utils"
	"math/big"

	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

	packedPrice, err := packUint256(r.Price.BigInt())
	if err != nil {
		return nil, err
	}

	return utils.AbiEncodeResponseTaskParam(r.TaskIndex, packedPrice)
}

func (p *ProcessRequestPriceFeedResultHandler) GetData(msg proto.Message) ([]byte, error) {
	return p.getAbiEncodeData(msg)
}

func (p *ProcessRequestPriceFeedResultHandler) GetDigest(msg proto.Message) ([]byte, error) {
	data, err := p.getAbiEncodeData(msg)
	if err != nil {
		return nil, err
	}
	return utils.DigestKeccak256(data), nil
}
