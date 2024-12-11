package result_handlers

import (
	"github.com/cosmos/gogoproto/proto"
	"intellix/pkg/utils"
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

	return utils.AbiEncodeResponseTaskParam(r.TaskIndex, r.ScriptOutData)
}

func (p *ProcessorRequestResHandler) GetDigest(msg proto.Message) ([]byte, error) {
	r, ok := msg.(*types.RequestScriptOut)
	if !ok {
		return nil, nil
	}
	data, err := utils.AbiEncodeResponseTaskParam(r.TaskIndex, r.ScriptOutData)
	if err != nil {
		return nil, err
	}
	return utils.DigestKeccak256(data), nil
}
