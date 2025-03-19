package server

import (
	"intellix/dvs/vrf/types"
)

type ResponseServer struct {
	Server
}

// NewRequestServer returns an implementation of the DvsProcessRequestServer interface for the provided Server.
func NewResponseServer(server Server) types.VRFMsgResponseServer {
	return &ResponseServer{
		Server: server,
	}
}

var _ types.VRFMsgResponseServer = ResponseServer{}
