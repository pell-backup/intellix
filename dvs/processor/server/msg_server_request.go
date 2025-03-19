package server

import (
	"bytes"
	context "context"
	"fmt"
	"time"

	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
	"github.com/IntelliXLabs/iwasm/api"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/x/authz"

	"intellix/dvs/processor/types"
	processortypes "intellix/x/processor/types"
)

var _ types.DVSRequestServer = Server{}

func (s Server) RequestScript(ctx context.Context, in *types.RequestScriptIn) (*types.RequestScriptOut, error) {
	pkgContext := sdktypes.UnwrapContext(ctx)
	s.logger.Info("RequestScript", "in", fmt.Sprintf("%+v", in))

	instance, runtime, scriptConfig, err := s.loadWasmScript(pkgContext, in.ScriptId)
	if err != nil {
		return nil, err
	}
	defer runtime.Dispose()
	defer instance.Dispose()

	data, err := s.fetchDataByExecWasmFetchingScript(pkgContext, instance, scriptConfig, in.ScriptParam)
	if err != nil {
		return nil, err
	}

	aggrDataIn, err := s.waitForEnoughOperateVote(pkgContext, in, data)
	if err != nil {
		return nil, err
	}

	aggrData, digest, err := s.aggrDataByExecWasmAggrScript(pkgContext, instance, scriptConfig, aggrDataIn)
	if err != nil {
		return nil, err
	}

	return &types.RequestScriptOut{
		TaskIndex:     in.TaskIndex,
		ScriptOutData: aggrData,
		DataDigest:    digest,
	}, nil
}

func (s Server) loadWasmScript(ctx sdktypes.Context, scriptId uint64) (api.InstanceResult, api.RuntimeResult, []byte, error) {
	conn := s.clientCtx.GRPCClient
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

func (s Server) fetchDataByExecWasmFetchingScript(ctx sdktypes.Context, instance api.InstanceResult, scriptConfig []byte, scriptParam []byte) ([]byte, error) {
	dataRes, err := instance.PrepareData(scriptConfig, scriptParam)
	if err != nil {
		return nil, err
	}
	data, err := dataRes.Data()
	dataRes.Dispose()

	return data, err
}

func (s Server) waitForEnoughOperateVote(ctx sdktypes.Context, in *types.RequestScriptIn, srcData []byte) ([][]byte, error) {
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
	if err := s.SignAndBroadcastTx(ctx, &voteIn); err != nil {
		s.logger.Error("waitForEnoughOperateVote SignAndBroadcastTx error: " + err.Error())
		return nil, fmt.Errorf("failed to broadcast VoteRequestProcessorIn for data error: %w", err)
	}

	// listen and collect [N-N+M] block
	var voteTxs []processortypes.MsgVoteRequestProcessor
	firstTxBlock := int64(0)
	for {
		block, err := s.GetLatestBlock(ctx)
		if err != nil {
			s.logger.Error("collectVoteRequestPriceFeed GetLatestBlock error: " + err.Error())
			return nil, fmt.Errorf("failed to get latest block: %w", err)
		}

		newTxs := s.processBlockTxs(ctx, block, in.RequestId)
		voteTxs = append(voteTxs, newTxs...)

		if firstTxBlock == 0 && len(newTxs) > 0 {
			firstTxBlock = block.Header.Height
		}

		// check N-N+M
		if s.shouldStopCollecting(ctx, firstTxBlock, block.Header.Height) {
			s.logger.Info("collectVoteRequestPriceFeed stop collecting", "block_height", block.Header.Height)
			break
		}

		// wait for next block
		time.Sleep(time.Millisecond * 10)
	}

	var resp [][]byte
	for _, v := range voteTxs {
		resp = append(resp, v.ScriptResp)
	}

	return resp, nil
}

func (s Server) processBlockTxs(ctx context.Context, block *cmttypes.Block, requestID []byte) []processortypes.MsgVoteRequestProcessor {
	var priceFeedTxs []processortypes.MsgVoteRequestProcessor
	//d.logger.Info("Processing block", "height", block.Header.Height, "tx_count", len(block.Data.Txs))

	for _, tx := range block.Data.Txs {
		if msg, ok := s.isVoteRequestTx(tx); ok && bytes.Equal(msg.RequestId, requestID) {
			priceFeedTxs = append(priceFeedTxs, *msg)
		}
	}

	return priceFeedTxs
}

func (s Server) isVoteRequestTx(tx cmttypes.Tx) (*processortypes.MsgVoteRequestProcessor, bool) {
	decoder := s.clientCtx.TxConfig.TxDecoder()
	data, err := decoder(tx)
	if err != nil {
		s.logger.Error("TxDecoder decode tx error: " + err.Error())
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
			s.logger.Error("Failed to get inner messages from authz", "error", err)
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

func (s Server) shouldStopCollecting(ctx sdktypes.Context, firstTxBlock, currentBlock int64) bool {
	if firstTxBlock == 0 {
		return false
	}
	return currentBlock >= firstTxBlock+s.waitBlockCount
}

func (s Server) aggrDataByExecWasmAggrScript(ctx sdktypes.Context, instance api.InstanceResult, scriptConfig []byte, datas [][]byte) ([]byte, []byte, error) {
	dataRes, err := instance.Aggregate(scriptConfig, datas, []byte("first"))
	if err != nil {
		return nil, nil, err
	}

	data, digest, err := dataRes.Data()
	dataRes.Dispose()

	return data, digest, err
}
