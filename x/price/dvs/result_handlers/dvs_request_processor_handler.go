package resulthandlers

import (
	"github.com/cosmos/gogoproto/proto"
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

	return r.DataDigest, nil
}
