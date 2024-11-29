package server

import (
	"intellix/x/price/dvs/types"
)

type RequestServer struct {
	Server
}

// NewDvsProcessRequestServer returns an implementation of the DvsProcessRequestServer interface
// for the provided Server.
func NewRequestServer(server Server) types.DvsProcessRequestServer {
	return &RequestServer{
		Server: server,
	}
}

var _ types.DvsProcessRequestServer = RequestServer{}
