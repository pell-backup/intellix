package handler

import (
	"fmt"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/crypto"
	"intellix/dvs/vrf/types"
	"math/big"
	"strings"
)

type VRFResultHandler struct {
}

func NewVRFResultHandler() *VRFResultHandler {
	return &VRFResultHandler{}
}

func (p *VRFResultHandler) GetData(msg proto.Message) ([]byte, error) {
	r, ok := msg.(*types.VRFTaskResponse)
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

func ParseData(data []byte) ([]*big.Int, error) {
	s := string(data)

	if s == "" {
		return nil, nil
	}

	parts := strings.Split(s, ",")
	result := make([]*big.Int, 0, len(parts))

	for _, part := range parts {
		bi, ok := new(big.Int).SetString(part, 10)
		if !ok {
			return nil, fmt.Errorf("parse big.Int failed for part: %s", part)
		}
		result = append(result, bi)
	}

	return result, nil
}
