package handler

import (
	"math/big"
	"strings"

	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/accounts/abi"

	"intellix/dvs"
	"intellix/dvs/vrf/types"
)

type VRFResultHandler struct {
}

func NewVRFResultHandler() *VRFResultHandler {
	return &VRFResultHandler{}
}

func (p *VRFResultHandler) getAbiEncodeData(msg proto.Message) ([]byte, error) {
	r, ok := msg.(*types.VRFTaskResponse)
	if !ok {
		return nil, nil
	}

	abiJSON := `
    [
      {
        "name": "foo",
        "type": "function",
        "inputs": [
          {
            "type": "uint256[]",
            "name": "randomWords"
          }
        ],
        "outputs": []
      }
    ]
    `
	myAbi, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, err
	}

	method := myAbi.Methods["foo"]
	randomWords := make([]*big.Int, len(r.RandomNumber))
	for i, v := range r.RandomNumber {
		randomWords[i] = v.BigInt()
	}

	packedArgsOnly, err := method.Inputs.Pack(randomWords)
	if err != nil {
		return nil, err
	}

	return dvs.AbiEncodeResponseTaskParam(r.TaskIndex, packedArgsOnly)
}

func (p *VRFResultHandler) GetData(msg proto.Message) ([]byte, error) {
	return p.getAbiEncodeData(msg)
}

func (p *VRFResultHandler) GetDigest(msg proto.Message) ([]byte, error) {
	data, err := p.getAbiEncodeData(msg)
	if err != nil {
		return nil, err
	}
	return dvs.DigestKeccak256(data), nil
}
