package handler

import (
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/crypto"
	"intellix/dvs/vrf/types"
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
	return []byte(r.RandomNumber.String()), nil
}

func (p *VRFResultHandler) GetDigest(msg proto.Message) ([]byte, error) {
	r, ok := msg.(*types.GenerateRandomNumberResponse)
	if !ok {
		return nil, nil
	}
	hashBytes := crypto.Keccak256(r.RandomNumber.BigInt().Bytes())
	return hashBytes, nil
}
