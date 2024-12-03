package server

import (
	"intellix/x/processor/dvs/types"

	"github.com/IntelliXLabs/iwasm/api"
)

type RequestServer struct {
	runtime api.RuntimeResult
}

// NewDvsProcessRequestServer returns an implementation of the DvsProcessRequestServer interface
// for the provided Server.
func NewRequestServer() types.DVSRequestServer {
	return &RequestServer{
		runtime: api.NewRuntime(),
	}
}

var _ types.DVSRequestServer = RequestServer{}
