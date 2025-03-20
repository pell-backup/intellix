package submitter

import (
	"context"
	"cosmossdk.io/math"
	"encoding/json"
	"fmt"
	contractdataoracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"intellix/gateway/types"
	"sync"
)

func (s *Submitter) RespondToPriceTask(req *types.RPCVoteFinalizedRequestIn, resp *types.RespondToTaskResponse) error {
	err := s.handlePriceTaskResponse(context.Background(), req)
	if err != nil {
		resp.Error = err.Error()
		return err
	}
	resp.Error = ""
	return nil
}

func (s *Submitter) handlePriceTaskResponse(ctx context.Context, response *types.RPCVoteFinalizedRequestIn) error {
	value, loaded := s.taskMap.LoadOrStore(response.TaskRaw.TaskIndex, response)
	if loaded {
		existingResponse := value.(*types.RPCVoteFinalizedRequestIn)
		if s.shouldReplacePriceTaskResponse(ctx, existingResponse, response) {
			s.taskMap.Store(response.TaskRaw.RequestID, response)
			return s.wrapPriceTaskSubmitToChain(ctx, response)
		}
	} else {
		return s.wrapPriceTaskSubmitToChain(ctx, response)
	}
	return nil
}

func (s *Submitter) shouldReplacePriceTaskResponse(ctx context.Context, existing, new *types.RPCVoteFinalizedRequestIn) bool {
	// TODO: add security threshold comparison
	return false
}

func (s *Submitter) wrapPriceTaskSubmitToChain(ctx context.Context, request *types.RPCVoteFinalizedRequestIn) error {
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	var err error
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.logger.Error("Failed to submit vote finalized request", "error", r)
				err = fmt.Errorf("%v", r)
			}
			wg.Done()
		}()
		err = s.submitPriceTaskResultToChain(ctx, request)
	}()
	wg.Wait()
	return err
}

func (s *Submitter) submitPriceTaskResultToChain(ctx context.Context, response *types.RPCVoteFinalizedRequestIn) error {
	chainConn, ok := s.chainConnections[uint64(response.ChainID)]
	if !ok {
		return fmt.Errorf("no connection found for chain ID %d", response.ChainID)
	}

	jsData, _ := json.Marshal(response)
	s.logger.Info("Submitter.submitPriceTaskResultToChain request", "data", string(jsData))

	// Validate BLS signature components
	if err := validateBLSComponents(response.ValidatedData); err != nil {
		s.logger.Error("Invalid BLS signature components", "error", err)
		return fmt.Errorf("invalid BLS components: %v", err)
	}

	feeTokenAddr, err := convertAddressToString(response.TaskRaw.FeeToken)
	if err != nil {
		s.logger.Error("Error converting fee token address", "err", err)
		return err
	}
	cbAddr, err := convertAddressToString(response.TaskRaw.CallbackAddress)
	if err != nil {
		s.logger.Error("Error converting callback address", "err", err)
		return err
	}

	paymentInt, ok := math.NewIntFromString(response.TaskRaw.Payment)
	if !ok {
		s.logger.Error("Error converting taskRaw payment", "payment", response.TaskRaw.Payment)
		return fmt.Errorf("error converting taskRaw payment")
	}
	task := contractdataoracle.IDataOracleTask{
		TaskType:                 math.NewInt(response.TaskRaw.TaskType).BigInt(),
		RequestId:                [32]byte(response.TaskRaw.RequestID),
		FeeToken:                 *feeTokenAddr,
		AdvanceDecode:            response.TaskRaw.AdvanceDecode,
		Payment:                  paymentInt.BigInt(),
		RequestData:              response.TaskRaw.RequestData,
		CallbackAddress:          *cbAddr,
		CallbackFunctionId:       [4]byte(response.TaskRaw.CallbackFunctionID),
		TaskCreatedBlock:         response.TaskRaw.TaskCreatedBlock,
		GroupNumbers:             response.TaskRaw.QuorumNumbers,
		GroupThresholdPercentage: response.TaskRaw.QuorumThresholdPercentage,
	}

	sign := contractdataoracle.IBLSSignatureVerifierNonSignerStakesAndSignature{
		NonSignerGroupBitmapIndices: response.ValidatedData.NonSignerQuorumBitmapIndices,
		NonSignerPubkeys:            convertNonSignersPubkeysG1(response.ValidatedData.NonSignersPubkeysG1),
		GroupApks:                   convertQuorumApks(response.ValidatedData.QuorumApksG1),
		ApkG2:                       convertApkG2(response.ValidatedData.SignersApkG2),
		Sigma:                       convertSigma(response.ValidatedData.SignersAggSigG1),
		GroupApkIndices:             response.ValidatedData.QuorumApkIndices,
		TotalStakeIndices:           response.ValidatedData.TotalStakeIndices,
		NonSignerStakeIndices:       response.ValidatedData.NonSignerStakeIndices,
	}

	taskResp := contractdataoracle.IDataOracleTaskResponse{
		ReferenceTaskIndex: response.TaskRaw.TaskIndex,
		Data:               response.RespToTaskData,
	}

	authOpts, err := s.getAuthOpts(response.ChainID)
	if err != nil {
		return err
	}

	// For debug purpose,
	// set gas limit to 1000000 to bypass transaction pre-execution and force broadcast
	// authOpts.GasLimit = 1000000

	if conf, ok := s.cfg.Chains[uint64(response.ChainID)]; ok {
		if conf.GasLimit > 0 {
			authOpts.GasLimit = conf.GasLimit
		}
	}

	transaction, err := chainConn.contractDataOracle.ResponseToTask(authOpts, task, taskResp, sign)
	if err != nil {
		s.logger.Error("Error assembling RequestPrice tx",
			"chainID", response.ChainID,
			"err", err,
			"task", fmt.Sprintf("%+v", task),
			"taskResp", fmt.Sprintf("%+v", taskResp),
			"sign", fmt.Sprintf("%+v", sign))
		return err
	}

	if transaction != nil {
		err = s.queryTransactionResult(ctx, transaction, chainConn.ethClient)
		if err != nil {
			return fmt.Errorf("chain %d: %v", response.ChainID, err)
		}
	}

	return nil
}
