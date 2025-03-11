package server

import (
	"intellix/dvs/price/types"
	"intellix/pkg/tx_listener"
	pricetypes "intellix/x/price/types"
)

type RequestServer struct {
	Server
	PriceListener       tx_listener.ChainListenerIFace[string, *pricetypes.MsgVoteRequestPriceFeed, *pricetypes.MsgVoteRequestPriceFeed]
	tickConverterConfig map[string]map[string]string
}

// NewRequestServer returns an implementation of the DvsProcessRequestServer interface
// for the provided Server.
func NewRequestServer(server Server) types.DVSRequestServer {
	s := &RequestServer{
		Server:              server,
		tickConverterConfig: server.tickConverterConfig,
	}

	s.PriceListener = tx_listener.NewChainListener(
		server.logger, server.clientCtx,
		server.wsEndpoint,
		"tm.event='Tx' AND eventType='finalized_price_feed'", 1000,
		s.PriceEventHandler, s.PriceBlockHandler,
	)

	s.PriceListener.Start()

	return s
}

var _ types.DVSRequestServer = &RequestServer{}
