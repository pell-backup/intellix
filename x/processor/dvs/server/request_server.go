package server

import (
	"bytes"
	context "context"
	"fmt"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	pkgcontext "intellix/pkg/context"
	"intellix/x/processor/dvs/types"
	processortypes "intellix/x/processor/types"
	"time"

	"github.com/IntelliXLabs/iwasm/api"
)

type RequestServer struct {
	Server
}

func NewRequestServer(s Server) types.DVSRequestServer {
	return &RequestServer{
		Server: s,
	}
}

var _ types.DVSRequestServer = RequestServer{}

func (r RequestServer) RequestScript(ctx context.Context, in *types.RequestScriptIn) (*types.RequestScriptOut, error) {
	pkgContext := pkgcontext.UnwrapContext(ctx)
	r.logger.Info("RequestScript", "in", fmt.Sprintf("%+v", in))

	instance, runtime, scriptConfig, err := r.loadWasmScript(pkgContext, in.ScriptId)
	if err != nil {
		return nil, err
	}
	defer runtime.Dispose()
	defer instance.Dispose()

	data, err := r.fetchDataByExecWasmFetchingScript(pkgContext, instance, scriptConfig, in.ScriptParam)
	if err != nil {
		return nil, err
	}

	aggrDataIn, err := r.waitForEnoughOperateVote(pkgContext, in, data)
	if err != nil {
		return nil, err
	}

	aggrData, digest, err := r.aggrDataByExecWasmAggrScript(pkgContext, instance, scriptConfig, aggrDataIn)
	if err != nil {
		return nil, err
	}

	return &types.RequestScriptOut{
		TaskIndex:     in.TaskIndex,
		ScriptOutData: aggrData,
		DataDigest:    digest,
	}, nil
}

func (r RequestServer) loadWasmScript(ctx pkgcontext.Context, scriptId uint64) (api.InstanceResult, api.RuntimeResult, []byte, error) {
	conn := r.clientCtx.GRPCClient
	queryClient := processortypes.NewQueryClient(conn)

	req := &processortypes.QueryShowProcessorRequest{
		Id: scriptId,
	}

	resp, err := queryClient.ShowProcessor(ctx, req)
	if err != nil {
		return api.InstanceResult{}, api.RuntimeResult{}, nil, fmt.Errorf("failed to query processor: %w", err)
	}

	runtime := api.NewRuntime()
	if runtime.Err() != nil {
		defer runtime.Dispose()
		return api.InstanceResult{}, api.RuntimeResult{}, nil, fmt.Errorf("failed to create runtime: %w", runtime.Err())
	}

	instance, err := runtime.CreateInstance(resp.Processor.WasmCode)
	if err != nil {
		return api.InstanceResult{}, api.RuntimeResult{}, nil, fmt.Errorf("failed to create instance: %w", err)
	}
	if instance.Err() != nil {
		return api.InstanceResult{}, api.RuntimeResult{}, nil, fmt.Errorf("failed to create instance: %w", instance.Err())
	}

	return instance, runtime, resp.Processor.Config, err
}

func (r RequestServer) fetchDataByExecWasmFetchingScript(ctx pkgcontext.Context, instance api.InstanceResult, scriptConfig []byte, scriptParam []byte) ([]byte, error) {
	dataRes, err := instance.PrepareData(scriptConfig, scriptParam)
	if err != nil {
		return nil, err
	}
	data, err := dataRes.Data()
	dataRes.Dispose()

	return data, err
}

func (r RequestServer) waitForEnoughOperateVote(ctx pkgcontext.Context, in *types.RequestScriptIn, srcData []byte) ([][]byte, error) {
	// broadcast VoteRequestScript
	voteIn := processortypes.MsgVoteRequestProcessor{
		TaskIndex:                 in.TaskIndex,
		RequestId:                 in.RequestId,
		FeeToken:                  in.FeeToken,
		Payment:                   in.Payment,
		RequestData:               in.RequestData,
		CallbackAddress:           in.CallbackAddress,
		CallbackFunctionId:        in.CallbackFunctionId,
		TaskCreatedBlock:          in.TaskCreatedBlock,
		QuorumNumbers:             in.QuorumNumbers,
		QuorumThresholdPercentage: in.QuorumThresholdPercentage,
		ScriptId:                  in.ScriptId,
		ScriptResp:                srcData,
	}
	if err := r.Server.SignAndBroadcastTx(ctx, &voteIn); err != nil {
		r.logger.Error("waitForEnoughOperateVote SignAndBroadcastTx error: " + err.Error())
		return nil, fmt.Errorf("failed to broadcast VoteRequestProcessorIn for data error: %w", err)
	}

	// listen and collect [N-N+M] block
	var voteTxs []processortypes.MsgVoteRequestProcessor
	firstTxBlock := int64(0)
	for {
		block, err := r.Server.GetLatestBlock(ctx)
		if err != nil {
			r.logger.Error("collectVoteRequestPriceFeed GetLatestBlock error: " + err.Error())
			return nil, fmt.Errorf("failed to get latest block: %w", err)
		}

		newTxs := r.processBlockTxs(ctx, block, in.RequestId)
		voteTxs = append(voteTxs, newTxs...)

		if firstTxBlock == 0 && len(newTxs) > 0 {
			firstTxBlock = block.Header.Height
		}

		// check N-N+M
		if r.shouldStopCollecting(ctx, firstTxBlock, block.Header.Height) {
			r.logger.Info("collectVoteRequestPriceFeed stop collecting", "block_height", block.Header.Height)
			break
		}

		// wait for next block
		ctx = ctx.WithBlockHeight(block.Header.Height + 1)
		time.Sleep(time.Millisecond * 10)
	}

	var resp [][]byte
	for _, v := range voteTxs {
		resp = append(resp, v.ScriptResp)
	}

	return resp, nil
}

func (r RequestServer) processBlockTxs(ctx context.Context, block *cmttypes.Block, requestID []byte) []processortypes.MsgVoteRequestProcessor {
	var priceFeedTxs []processortypes.MsgVoteRequestProcessor
	//d.logger.Info("Processing block", "height", block.Header.Height, "tx_count", len(block.Data.Txs))

	for _, tx := range block.Data.Txs {
		if msg, ok := r.isVoteRequestTx(tx); ok && bytes.Equal(msg.RequestId, requestID) {
			priceFeedTxs = append(priceFeedTxs, *msg)
		}
	}

	return priceFeedTxs
}

func (r RequestServer) isVoteRequestTx(tx cmttypes.Tx) (*processortypes.MsgVoteRequestProcessor, bool) {
	decoder := r.clientCtx.TxConfig.TxDecoder()
	data, err := decoder(tx)
	if err != nil {
		r.logger.Error("TxDecoder decode tx error: " + err.Error())
		return nil, false
	}

	msgs := data.GetMsgs()
	if len(msgs) == 0 {
		return nil, false
	}

	msg := msgs[0]

	// Try to handle authz message
	if authzMsg, ok := msg.(*authz.MsgExec); ok {
		// Get the inner messages from authz
		innerMsgs, err := authzMsg.GetMessages()
		if err != nil {
			r.logger.Error("Failed to get inner messages from authz", "error", err)
			return nil, false
		}
		if len(innerMsgs) == 0 {
			return nil, false
		}
		// Use the first inner message
		msg = innerMsgs[0]
	}

	voteMsg, ok := msg.(*processortypes.MsgVoteRequestProcessor)
	if !ok {
		//d.logger.Error("msg is not VoteRequestProcessorIn", "msg", fmt.Sprintf("%+v", msg))
		return nil, false
	}

	return voteMsg, true
}

func (r RequestServer) shouldStopCollecting(ctx pkgcontext.Context, firstTxBlock, currentBlock int64) bool {
	if firstTxBlock == 0 {
		return false
	}
	return currentBlock >= firstTxBlock+r.waitBlockCount
}

func (r RequestServer) aggrDataByExecWasmAggrScript(ctx pkgcontext.Context, instance api.InstanceResult, scriptConfig []byte, datas [][]byte) ([]byte, []byte, error) {
	dataRes, err := instance.Aggregate(scriptConfig, datas, []byte("first"))
	if err != nil {
		return nil, nil, err
	}

	data, digest, err := dataRes.Data()
	dataRes.Dispose()

	return data, digest, err
}
