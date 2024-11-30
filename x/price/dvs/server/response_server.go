package server

import (
	"intellix/x/price/dvs/types"
)

type ResponseServer struct {
	Server
}

func NewResponseServer(server Server) types.DvsPostProcessRequestServer {
	return &ResponseServer{Server: server}
}

var _ types.DvsPostProcessRequestServer = ResponseServer{}
