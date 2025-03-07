package server

import (
	"intellix/dvs/price/types"
)

type RequestServer struct {
	Server
	tickConverterConfig map[string]map[string]string
}

// NewDvsProcessRequestServer returns an implementation of the DvsProcessRequestServer interface
// for the provided Server.
func NewRequestServer(server Server) types.DVSRequestServer {
	return &RequestServer{
		Server:              server,
		tickConverterConfig: server.tickConverterConfig,
	}
}

var _ types.DVSRequestServer = RequestServer{}
