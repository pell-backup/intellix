package server

import (
	"intellix/dvs/vrf/types"
)

type RequestServer struct {
	Server
}

// NewDvsProcessRequestServer returns an implementation of the DvsProcessRequestServer interface
// for the provided Server.
func NewRequestServer(server Server) types.VRFMsgServerServer {
	return &RequestServer{
		Server: server,
	}
}

var _ types.VRFMsgServerServer = RequestServer{}
