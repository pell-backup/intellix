package server

import (
	"context"
	"fmt"
	"intellix/x/price/types"
	pricetypes "intellix/x/price/types"

	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

func (d *RequestServer) PriceEventHandler(ctx context.Context, event abci.Event) (string, *pricetypes.MsgVoteRequestPriceFeed, error) {
	var (
		operatorId, requestId string
	)

	if event.Type == "price.VoteRequestPriceFeed" {
		for _, attr := range event.Attributes {
			switch attr.Key {
			case "operator_id":
				operatorId = attr.Value
			case "request_id":
				requestId = attr.Value
			}
		}
	}

	priceFeed, err := d.queryVoteRequestPriceFeed(ctx, []byte(requestId), operatorId)
	if err != nil {
		return "", nil, err
	}

	d.logger.Info("receive price feed event", "price feed", fmt.Sprintf("%+v", priceFeed))

	return requestId, priceFeed, nil
}

func (d *RequestServer) queryVoteRequestPriceFeed(ctx context.Context, requestId []byte, operatorId string) (*pricetypes.MsgVoteRequestPriceFeed, error) {
	conn := d.clientCtx.GRPCClient
	queryClient := pricetypes.NewQueryClient(conn)

	req := &types.QueryVoteRequestPriceFeedReq{
		RequestId:  requestId,
		OperatorId: operatorId,
	}

	resp, err := queryClient.QueryVoteRequestPriceFeed(ctx, req)
	if err != nil {
		return nil, err
	}
	var price []*pricetypes.VoteRequestPriceFeed
	for _, p := range resp.Price {
		price = append(price, &pricetypes.VoteRequestPriceFeed{
			Price:  p.Price,
			Source: p.Source,
		})
	}

	return &pricetypes.MsgVoteRequestPriceFeed{
		TaskIndex:    resp.TaskIndex,
		OperatorId:   resp.OperatorId,
		RequestId:    resp.RequestId,
		BaseSymbol:   resp.BaseSymbol,
		QuoteSymbol:  resp.QuoteSymbol,
		Price:        price,
		Timestamp:    resp.Timestamp,
		BlockHeight:  resp.BlockHeight,
		Sender:       resp.Sender,
		BlsSignature: resp.BlsSignature,
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
