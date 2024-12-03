package server

import (
	context "context"
	pkgcontext "intellix/pkg/context"
	"intellix/x/processor/dvs/types"

	"github.com/IntelliXLabs/iwasm/api"
)

type RequestServer struct {
	runtime api.RuntimeResult
}

func NewRequestServer() types.DVSRequestServer {
	return &RequestServer{
		runtime: api.NewRuntime(),
	}
}

var _ types.DVSRequestServer = RequestServer{}

func (r RequestServer) RequestScript(ctx context.Context, in *types.RequestScriptIn) (*types.RequestScriptOut, error) {
	//TODO implement me
	panic("implement me")
}

func (r RequestServer) loadWasmScript(ctx pkgcontext.Context, scriptId uint64) ([]byte, error) {
	// TODO: query wasm script from abci
	return nil, nil
}

func (r RequestServer) fetchDataByExecWasmFetchingScript(ctx pkgcontext.Context, scriptBytes []byte) ([]byte, error) {
	// TODO: execute wasm script to fetch data
	return nil, nil
}

func (r RequestServer) waitForEnoughOperate(ctx pkgcontext.Context) ([][]byte, error) {
	return nil, nil
}

func (r RequestServer) aggrDataByExecWasmAggrScript(ctx pkgcontext.Context, datas [][]byte) ([]byte, error) {
	return nil, nil
}
