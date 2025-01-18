package server

import (
	"intellix/dvs/price/types"
	"intellix/pkg/tx_listener"
	pricetypes "intellix/x/price/types"
)

type RequestServer struct {
	Server
	PriceListener tx_listener.ChainListenerIFace[string, *pricetypes.MsgVoteRequestPriceFeed, *pricetypes.MsgVoteRequestPriceFeed]
}

// NewDvsProcessRequestServer returns an implementation of the DvsProcessRequestServer interface
// for the provided Server.
func NewRequestServer(server Server) types.DVSRequestServer {
	s := &RequestServer{
		Server: server,
	}

	s.PriceListener = tx_listener.NewChainListener(
		server.logger, server.clientCtx,
		server.wsEndpoint,
		"tm.event='Tx' AND message.action='VoteRequestPriceFeed'", 1000,
		s.PriceEventHandler, s.PriceBlockHandler,
	)

	s.PriceListener.Start()

	return s
}

var _ types.DVSRequestServer = &RequestServer{}
