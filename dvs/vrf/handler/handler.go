package handler

import (
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/crypto"
	"intellix/dvs/vrf/types"
	"strings"
)

type VRFResultHandler struct {
}

func NewVRFResultHandler() *VRFResultHandler {
	return &VRFResultHandler{}
}

func (p *VRFResultHandler) GetData(msg proto.Message) ([]byte, error) {
	r, ok := msg.(*types.GenerateRandomNumberResponse)
	if !ok {
		return nil, nil
	}

	var sb strings.Builder
	for i, num := range r.RandomNumber {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(num.String())
	}
	return []byte(sb.String()), nil
}

func (p *VRFResultHandler) GetDigest(msg proto.Message) ([]byte, error) {
	data, err := p.GetData(msg)
	if err != nil {
		return nil, err
	}

	hashBytes := crypto.Keccak256(data)
	return hashBytes, nil
}
