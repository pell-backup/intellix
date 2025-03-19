package result_handlers

import (
	"github.com/cosmos/gogoproto/proto"

	"intellix/dvs"
	"intellix/dvs/processor/types"
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

	return dvs.AbiEncodeResponseTaskParam(r.TaskIndex, r.ScriptOutData)
}

func (p *ProcessorRequestResHandler) GetDigest(msg proto.Message) ([]byte, error) {
	r, ok := msg.(*types.RequestScriptOut)
	if !ok {
		return nil, nil
	}
	data, err := dvs.AbiEncodeResponseTaskParam(r.TaskIndex, r.ScriptOutData)
	if err != nil {
		return nil, err
	}
	return dvs.DigestKeccak256(data), nil
}
