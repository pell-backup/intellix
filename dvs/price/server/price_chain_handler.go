package server

import (
	"context"
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	pricetypes "intellix/x/price/types"
	"strconv"
)

func (d *RequestServer) PriceEventHandler(ctx context.Context, event abci.Event) (string, *pricetypes.MsgVoteRequestPriceFeed, error) {
	var (
		taskIndex, operatorId, requestId string
	)

	if event.Type == "price.VoteRequestPriceFeed" {
		for _, attr := range event.Attributes {
			switch attr.Key {
			case "task_index":
				taskIndex = attr.Value
			case "operator_id":
				operatorId = attr.Value
			case "request_id":
				requestId = attr.Value
			}
		}
	}
	// TODO: send cosmos query tx

	taskIndexInt, err := strconv.ParseInt(taskIndex, 10, 32)
	if err != nil {
		return "", nil, err
	}
	return requestId, &pricetypes.MsgVoteRequestPriceFeed{
		TaskIndex:  uint32(taskIndexInt),
		OperatorId: operatorId,
		RequestId:  []byte(requestId),
	}, nil
}

func (d *RequestServer) isVoteRequestPriceFeedTx(tx cmttypes.Tx) (*pricetypes.MsgVoteRequestPriceFeed, bool) {
	// check tx is VoteRequestPriceFeed
	decoder := d.clientCtx.TxConfig.TxDecoder()
	data, err := decoder(tx)
	if err != nil {
		d.logger.Error("TxDecoder decode tx error: " + err.Error())
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
			d.logger.Error("Failed to get inner messages from authz", "error", err)
			return nil, false
		}
		if len(innerMsgs) == 0 {
			return nil, false
		}
		// Use the first inner message
		msg = innerMsgs[0]
	}

	voteMsg, ok := msg.(*pricetypes.MsgVoteRequestPriceFeed)
	if !ok {
		//d.logger.Error("msg is not MsgVoteRequestPriceFeed", "msg", fmt.Sprintf("%+v", msg))
		return nil, false
	}

	return voteMsg, true
}

func (d *RequestServer) PriceBlockHandler(ctx context.Context, block *cmttypes.Block) (string, *pricetypes.MsgVoteRequestPriceFeed, error) {
	for _, tx := range block.Data.Txs {
		// only collect VoteRequestPriceFeed
		//d.logger.Info("Processing transaction", "tx_hash", tx.Hash())
		if msg, ok := d.isVoteRequestPriceFeedTx(tx); ok {
			return string(msg.RequestId), &pricetypes.MsgVoteRequestPriceFeed{
				TaskIndex:   msg.TaskIndex,
				OperatorId:  msg.OperatorId,
				RequestId:   msg.RequestId,
				BaseSymbol:  msg.BaseSymbol,
				QuoteSymbol: msg.QuoteSymbol,
				Price:       msg.Price,
				Timestamp:   msg.Timestamp,
				BlockHeight: uint64(block.Header.Height),
			}, nil
		}
	}

	return "", nil, fmt.Errorf("no vote request price feed tx")
}
