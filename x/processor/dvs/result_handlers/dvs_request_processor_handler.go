package result_handlers

import (
	"github.com/cosmos/gogoproto/proto"
	"golang.org/x/crypto/sha3"
	"intellix/x/processor/dvs/types"
)

type ProcessorRequestResHandler struct {
}

func NewProcessorRequestResHandler() *ProcessorRequestResHandler {
	return &ProcessorRequestResHandler{}
}

func (p *ProcessorRequestResHandler) GetData(msg proto.Message) ([]byte, error) {
	r, ok := msg.(*types.RequestScriptOut)
	if !ok {
		return nil, nil
	}

	return r.ScriptOutData, nil
}

func (p *ProcessorRequestResHandler) GetDigest(msg proto.Message) ([]byte, error) {
	r, ok := msg.(*types.RequestScriptOut)
	if !ok {
		return nil, nil
	}

	//return r.DataDigest, nil
	// XXX: calc digest by script out
	var taskResponseDigest [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(r.ScriptOutData)
	copy(taskResponseDigest[:], hasher.Sum(nil)[:32])

	return taskResponseDigest[:], nil
}
