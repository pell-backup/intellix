package server

import (
	context "context"
	processortypes "intellix/x/processor/types"

	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

func (r *RequestServer) ProcessorEventHandler(ctx context.Context, event abci.Event) (string, *processortypes.MsgVoteRequestProcessor, error) {

	return "", nil, nil
}

func (r *RequestServer) ProcessorBlockHandler(ctx context.Context, block *cmttypes.Block) (string, *processortypes.MsgVoteRequestProcessor, error) {
	for _, tx := range block.Data.Txs {
		if msg, ok := r.isVoteRequestTx(tx); ok {
			return string(msg.RequestId), msg, nil
		}
	}
	return "", nil, nil
}

func (r *RequestServer) isVoteRequestTx(tx cmttypes.Tx) (*processortypes.MsgVoteRequestProcessor, bool) {
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
