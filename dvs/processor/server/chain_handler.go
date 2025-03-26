package server

import (
	context "context"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/x/authz"

	processortypes "intellix/x/processor/types"
)

func (r Server) ProcessorEventHandler(ctx context.Context, event abci.Event) (string, *processortypes.MsgVoteRequestProcessor, error) {
	var requestId string
	var operatorId string

	for _, attr := range event.Attributes {
		switch attr.Key {
		case "request_id":
			requestId = attr.Value
		case "operator_id":
			operatorId = attr.Value
		}
	}

	processor, err := r.queryVoteRequestProcessor(ctx, []byte(requestId), operatorId)
	if err != nil {
		return "", nil, err
	}

	r.logger.Info("receive processor event", "processor", fmt.Sprintf("%+v", processor))

	return requestId, processor, nil
}

func (r Server) queryVoteRequestProcessor(ctx context.Context, requestId []byte, operatorId string) (*processortypes.MsgVoteRequestProcessor, error) {
	conn := r.clientCtx.GRPCClient
	queryClient := processortypes.NewQueryClient(conn)

	req := &processortypes.QueryVoteRequestProcessorRequest{
		RequestId:  requestId,
		OperatorId: operatorId,
	}

	resp, err := queryClient.QueryVoteRequestProcessor(ctx, req)
	if err != nil {
		return nil, err
	}

	return &processortypes.MsgVoteRequestProcessor{
		TaskIndex:                 resp.TaskIndex,
		RequestId:                 resp.RequestId,
		FeeToken:                  resp.FeeToken,
		Payment:                   resp.Payment,
		RequestData:               resp.RequestData,
		CallbackAddress:           resp.CallbackAddress,
		CallbackFunctionId:        resp.CallbackFunctionId,
		TaskCreatedBlock:          resp.TaskCreatedBlock,
		QuorumNumbers:             resp.QuorumNumbers,
		QuorumThresholdPercentage: resp.QuorumThresholdPercentage,
		ScriptId:                  resp.ScriptId,
		ScriptResp:                resp.ScriptResp,
		Sender:                    resp.Sender,
		OperatorId:                resp.OperatorId,
		BlsSignature:              resp.BlsSignature,
	}, nil
}

func (r Server) ProcessorBlockHandler(ctx context.Context, block *cmttypes.Block) (string, *processortypes.MsgVoteRequestProcessor, error) {
	for _, tx := range block.Data.Txs {
		if msg, ok := r.isVoteRequestTx(tx); ok {
			return string(msg.RequestId), msg, nil
		}
	}
	return "", nil, nil
}

func (d Server) ProcessorMempoolEventHandler(ctx context.Context, msg *processortypes.MsgVoteRequestProcessor) (string, *processortypes.MsgVoteRequestProcessor, error) {
	return string(msg.RequestId), msg, nil
}

func (r Server) isVoteRequestTx(tx cmttypes.Tx) (*processortypes.MsgVoteRequestProcessor, bool) {
	decoder := r.clientCtx.TxConfig.TxDecoder()
	data, err := decoder(tx)
	if err != nil {
		r.logger.Error("TxDecoder decode tx error: " + err.Error())
		return nil, false
	}

	msgs := data.GetMsgs()
	if len(msgs) == 0 {
		return nil, false
	}

	msg := msgs[0]

	// Try to handle authz message
	if authzMsg, ok := msg.(*authz.MsgExec); ok {
		// Get the inner messages from authz
		innerMsgs, err := authzMsg.GetMessages()
		if err != nil {
			r.logger.Error("Failed to get inner messages from authz", "error", err)
			return nil, false
		}
		if len(innerMsgs) == 0 {
			return nil, false
		}
		// Use the first inner message
		msg = innerMsgs[0]
	}

	voteMsg, ok := msg.(*processortypes.MsgVoteRequestProcessor)
	if !ok {
		//d.logger.Error("msg is not VoteRequestProcessorIn", "msg", fmt.Sprintf("%+v", msg))
		return nil, false
	}

	return voteMsg, true
}
