package server

import (
	"bytes"
	context "context"
	"fmt"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"sort"
	"sync"
	"time"

	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
	"github.com/IntelliXLabs/iwasm/api"
	"github.com/cosmos/gogoproto/proto"

	"intellix/dvs/processor/types"
	"intellix/pkg/tx_listener"
	sdktypes "intellix/sdk/types"
	"intellix/sdk/utils"
	processortypes "intellix/x/processor/types"
)

type RequestServer struct {
	Server
	ProcessorListener tx_listener.ChainListenerIFace[string, *processortypes.MsgVoteRequestProcessor, *processortypes.MsgVoteRequestProcessor]
}

func NewRequestServer(s Server) types.DVSRequestServer {
	srv := &RequestServer{
		Server: s,
	}

	srv.ProcessorListener = tx_listener.NewChainListener(
		s.logger, s.clientCtx,
		s.wsEndpoint,
		"tm.event='Tx' AND eventType='vote_request_processor'", 1000,
		srv.ProcessorEventHandler, srv.ProcessorBlockHandler,
	)

	srv.ProcessorListener.Start()

	return srv
}

var _ types.DVSRequestServer = &RequestServer{}

func (r *RequestServer) RequestScript(ctx context.Context, in *types.RequestScriptIn) (*types.RequestScriptOut, error) {
	pkgContext := sdktypes.UnwrapContext(ctx)
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


func (r *RequestServer) loadWasmScript(ctx sdktypes.Context, scriptId uint64) (api.InstanceResult, api.RuntimeResult, []byte, error) {
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

func (r *RequestServer) waitForEnoughOperateVote(ctx sdktypes.Context, in *types.RequestScriptIn, srcData []byte) ([][]byte, error) {
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

	// sign
	sign := utils.SignWithBLS(r.Server.blsKeyPair, r.getMsgBytes(&voteIn))
	if sign == nil {
		return nil, fmt.Errorf("failed to sign VoteRequestProcessorIn")
	}
	voteIn.BlsSignature = sign

	if err := r.Server.SignAndBroadcastTx(ctx, &voteIn); err != nil {
		r.logger.Error("waitForEnoughOperateVote SignAndBroadcastTx error: " + err.Error())
		return nil, fmt.Errorf("failed to broadcast VoteRequestProcessorIn for data error: %w", err)
	}

	// wait for enough events or blocks
	events, err := r.collectVoteRequestProcessor(ctx, in)
	if err != nil {
		r.logger.Error("waitForEnoughOperateVote collectVoteRequestProcessor error: " + err.Error())
		return nil, fmt.Errorf("failed to collect VoteRequestProcessor: %w", err)
	}
	if len(events) == 0 {
		r.logger.Error("waitForEnoughOperateVote no enough events")
		return nil, fmt.Errorf("no enough events")
	}

	var resp [][]byte
	for _, v := range events {
		resp = append(resp, v.ScriptResp)
	}

	return resp, nil
}


func (r *RequestServer) collectEvents(ctx context.Context, operatorMaps map[string]*avsitypes.Operator, eventCh *tx_listener.EventChannel[*processortypes.MsgVoteRequestProcessor], eventData *[]*processortypes.MsgVoteRequestProcessor, waitChan chan struct{}) {
	eventDataByOperatorId := make(map[string]*processortypes.MsgVoteRequestProcessor)
	for _, data := range *eventData {
		eventDataByOperatorId[data.OperatorId] = data
	}

	mu := sync.Mutex{}
	select {
	case <-ctx.Done():
		return
	case data, ok := <-eventCh.Data:
		if !ok {
			return
		}
		// verify operator sign
		operator, ok := operatorMaps[data.Data.OperatorId]
		if ok && utils.VerifyBLSSignature(operator.Pubkeys, r.getMsgBytes(data.Data), data.Data.BlsSignature) == nil {
			mu.Lock()
			eventDataByOperatorId[data.Data.OperatorId] = data.Data
			mu.Unlock()
		}
		if len(eventDataByOperatorId) >= len(operatorMaps) {
			mu.Lock()
			// distinct by operator
			eventData = &[]*processortypes.MsgVoteRequestProcessor{}
			for _, data := range eventDataByOperatorId {
				*eventData = append(*eventData, data)
			}
			mu.Unlock()
			waitChan <- struct{}{}
			return
		}
	}
}

func (s Server) isVoteRequestTx(tx cmttypes.Tx) (*processortypes.MsgVoteRequestProcessor, bool) {
	decoder := s.clientCtx.TxConfig.TxDecoder()
	data, err := decoder(tx)
	if err != nil {
		s.logger.Error("TxDecoder decode tx error: " + err.Error())
		return nil, false
	}

	// subscribe events and blocks

	var (
		eventWaitChan      = make(chan struct{})
		blockWaitChan      = make(chan struct{})
		eventCh            = r.ProcessorListener.SubscribeEvents(1000)
		blockCh            = r.ProcessorListener.SubscribeBlocks(1000)
		operatorMaps       = make(map[string]*avsitypes.Operator)
		collectCtx, cancel = context.WithTimeout(ctx, 10*time.Minute)
	)
	for _, v := range ctx.Operators() {
		operatorMaps[string(v.Id)] = v
	}
	defer func() {
		cancel()

		r.ProcessorListener.UnsubscribeEvents(eventCh)
		r.ProcessorListener.UnsubscribeBlocks(blockCh)

		close(eventWaitChan)
		close(blockWaitChan)

		r.ProcessorListener.ClearBlocksByKey(key)
		r.ProcessorListener.ClearEventsByKey(key)
	}()

	go r.collectEvents(collectCtx, operatorMaps, eventCh, &eventData, eventWaitChan) // listen events
	go r.collectBlocks(collectCtx, blockHeight, blockCh, &blockData, blockWaitChan)  // listen blocks

	select {
	case <-blockWaitChan:
		return blockData, nil
	case <-eventWaitChan:
		return eventData, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("context done")
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
}

func (r *RequestServer) collectEvents(ctx context.Context, operatorMaps map[string]*avsitypes.Operator, eventCh *tx_listener.EventChannel[*processortypes.MsgVoteRequestProcessor], eventData *[]*processortypes.MsgVoteRequestProcessor, waitChan chan struct{}) {
	eventDataByOperatorId := make(map[string]*processortypes.MsgVoteRequestProcessor)
	for _, data := range *eventData {
		eventDataByOperatorId[data.OperatorId] = data
	}

	mu := sync.Mutex{}
	select {
	case <-ctx.Done():
		return
	case data, ok := <-eventCh.Data:
		if !ok {
			return
		}
		// verify operator sign
		operator, ok := operatorMaps[data.Data.OperatorId]
		if ok && utils.VerifyBLSSignature(operator.Pubkeys, r.getMsgBytes(data.Data), data.Data.BlsSignature) == nil {
			mu.Lock()
			eventDataByOperatorId[data.Data.OperatorId] = data.Data
			mu.Unlock()
		}
		if len(eventDataByOperatorId) >= len(operatorMaps) {
			mu.Lock()
			// distinct by operator
			eventData = &[]*processortypes.MsgVoteRequestProcessor{}
			for _, data := range eventDataByOperatorId {
				*eventData = append(*eventData, data)
			}
			mu.Unlock()
			waitChan <- struct{}{}
			return
		}
	}
}

func (r *RequestServer) collectBlocks(ctx context.Context, blockHeight int64, ch *tx_listener.BlockChannel[*processortypes.MsgVoteRequestProcessor], blockData *[]*processortypes.MsgVoteRequestProcessor, waitChan chan struct{}) {
	var mu = sync.Mutex{}
	select {
	case <-ctx.Done():
		return
	case <-ch.Done:
		return
	case data, ok := <-ch.Data:
		if !ok {
			return
		}
		if data.Height > blockHeight {
			mu.Lock()
			*blockData = append(*blockData, data.Data)
			mu.Unlock()
		}
		if data.Height >= blockHeight+r.waitBlockCount {
			waitChan <- struct{}{}
			return
		}
	}
}

func (r *RequestServer) getMsgBytes(msg *processortypes.MsgVoteRequestProcessor) []byte {
	b := &processortypes.MsgVoteRequestProcessor{
		TaskIndex:                 msg.TaskIndex,
		RequestId:                 msg.RequestId,
		FeeToken:                  msg.FeeToken,
		Payment:                   msg.Payment,
		RequestData:               msg.RequestData,
		CallbackAddress:           msg.CallbackAddress,
		CallbackFunctionId:        msg.CallbackFunctionId,
		TaskCreatedBlock:          msg.TaskCreatedBlock,
		QuorumNumbers:             msg.QuorumNumbers,
		QuorumThresholdPercentage: msg.QuorumThresholdPercentage,
		ScriptId:                  msg.ScriptId,
		ScriptResp:                msg.ScriptResp,
		Sender:                    msg.Sender,
		OperatorId:                msg.OperatorId,
	}

	bytes, _ := proto.Marshal(b)
	return bytes
}

func (r *RequestServer) verifyOperatorEvents(ctx sdktypes.Context, events []tx_listener.EventData[*processortypes.MsgVoteRequestProcessor]) ([]*processortypes.MsgVoteRequestProcessor, bool) {

	var operatorMaps = make(map[string]*avsitypes.Operator)
	var operatorVerifyMaps = make(map[string]bool)
	for _, v := range ctx.Operators() {
		operatorMaps[string(v.Id)] = v
	}

	var out []*processortypes.MsgVoteRequestProcessor
	for _, v := range events {
		if operator, ok := operatorMaps[v.Data.OperatorId]; ok {
			// verify operator's key
			err := utils.VerifyBLSSignature(operator.Pubkeys, r.getMsgBytes(v.Data), v.Data.BlsSignature)
			if err != nil {
				continue
			}
			operatorVerifyMaps[string(operator.Id)] = true
			out = append(out, v.Data)
		}
	}

	return out, len(operatorVerifyMaps) == len(operatorMaps)
}

func (r *RequestServer) checkAndChooseEnoughBlocks(ctx context.Context, blockHeight int64, blockDatas []tx_listener.BlockData[*processortypes.MsgVoteRequestProcessor]) ([]*processortypes.MsgVoteRequestProcessor, bool) {
	if len(blockDatas) == 0 {
		return nil, false
	}

	var out []*processortypes.MsgVoteRequestProcessor
	sort.Slice(blockDatas, func(i, j int) bool {
		return blockDatas[i].Height > blockDatas[j].Height
	})

	// choose the latest block
	maxHeightBlocks := blockDatas[len(blockDatas)-1].Height

	for _, data := range blockDatas {
		if data.Height >= blockHeight {
			out = append(out, data.Data)
		}
	}

	return out, maxHeightBlocks >= blockHeight+r.waitBlockCount
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
