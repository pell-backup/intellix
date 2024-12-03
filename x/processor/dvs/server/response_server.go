package server

import (
	"context"
	"github.com/IntelliXLabs/iwasm/api"
	dvstypes "intellix/pkg/pelldvs/types"
	"intellix/x/processor/dvs/types"
)

type ResponseServer struct {
	runtime api.RuntimeResult
}

func NewResponseServer() types.DVSResponseServer {
	return &ResponseServer{
		runtime: api.NewRuntime(),
	}
}

var _ types.DVSResponseServer = ResponseServer{}

func (r ResponseServer) ResponseScript(ctx context.Context, in *types.RequestScriptIn) (*types.ResponseScriptOut, error) {
	//TODO implement me
	panic("implement me")
}

func (r ResponseServer) voteData(ctx context.Context, in *types.RequestScriptIn, data []byte) error {
	return nil
}

func (r ResponseServer) responseToTask(ctx context.Context, in *types.RequestScriptIn, data []byte, validatedData *dvstypes.RequestPostRequestValidatedData) error {
	return nil
}
