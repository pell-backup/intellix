package dispatcher

import (
	"context"
	"cosmossdk.io/math"
	"encoding/hex"
	"errors"
	"fmt"
	interactortypes "github.com/0xPellNetwork/pelldvs-interactor/types"
	"github.com/0xPellNetwork/pelldvs/rpc/client/http"
	contractdataoracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ontio/ontology-crypto/keypair"
	"intellix/common"
	"intellix/dvs"
	pricetypes "intellix/dvs/price/types"
	processortypes "intellix/dvs/processor/types"
	"intellix/dvs/vrf/types"
	"strings"
)

func (d *Dispatcher) listenForNewTasks(chain *chainWatcher) {
	d.logger.Info("Starting task listener", "chainID", chain.chainID)

	newTaskChan := make(chan *contractdataoracle.ContractDataOracleNewTaskCreated)
	sub, err := chain.dataOracleContract.WatchNewTaskCreated(&bind.WatchOpts{}, newTaskChan, nil)
	if err != nil {
		d.logger.Error("Failed to watch for new tasks", "chainID", chain.chainID, "error", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			d.logger.Error("Subscription error", "chainID", chain.chainID, "error", err)
			return
		case newTask := <-newTaskChan:
			d.handleNewTask(chain.chainID, newTask)
		case <-d.Quit():
			return
		}
	}
}

func (d *Dispatcher) handleNewTask(chainID uint64, newTask *contractdataoracle.ContractDataOracleNewTaskCreated) {
	d.logger.Info("New task created",
		"chainID", chainID,
		"TaskIndex", newTask.TaskIndex,
		"RequestId", newTask.Task.RequestId,
		"TaskType", newTask.Task.TaskType,
	)

	taskData, err := d.serializeTaskByType(chainID, newTask)
	if err != nil {
		d.logger.Error("Failed to serialize task", "chainID", chainID, "error", err)
		return
	}

	err = d.dispatchTask(taskData, newTask.Task.GroupNumbers, chainID, uint32(newTask.Raw.BlockNumber), newTask.Task.GroupThresholdPercentage)
	if err != nil {
		d.logger.Error("Failed to dispatch task", "chainID", chainID, "error", err)
	}
}

func (d *Dispatcher) serializeTaskByType(chainID uint64, task *contractdataoracle.ContractDataOracleNewTaskCreated) ([]byte, error) {
	d.logger.Info("Serializing task by type",
		"chainID", chainID,
		"taskIndex", task.TaskIndex,
		"requestId", task.Task.RequestId,
		"taskType", task.Task.TaskType.Int64(),
	)

	var taskRequest sdk.Msg
	var err error

	switch task.Task.TaskType.Int64() {
	case dvs.TaskTypePriceFeed:
		taskRequest, err = d.serializePriceTask(chainID, task)
	case dvs.TaskTypeProcessor:
		taskRequest, err = d.serializeScriptTask(chainID, task)
	case dvs.TaskTypeVRFRandomNumber:
		taskRequest, err = d.serializeVRFTask(chainID, task)
	default:
		return nil, fmt.Errorf("invalid task type: %d", task.Task.TaskType.Int64())
	}

	if err != nil {
		return nil, err
	}
	if taskRequest == nil {
		return nil, errors.New("task request is nil")
	}

	return d.msgEncoder.EncodeMsgs(taskRequest)
}

func (d *Dispatcher) serializePriceTask(chainID uint64, newTask *contractdataoracle.ContractDataOracleNewTaskCreated) (sdk.Msg, error) {
	task := newTask.Task
	priceFeed, err := ParsePriceFeed(task.RequestData)
	if err != nil {
		d.logger.Error("Failed to parse price feed", "chainID", chainID, "error", err)
		return nil, err
	}
	taskRequest := &pricetypes.RequestPriceFeedIn{
		Task: &pricetypes.TaskRequest{
			TaskIndex:                 newTask.TaskIndex,
			RequestId:                 task.RequestId[:],
			FeeToken:                  task.FeeToken.Hex(),
			Payment:                   math.NewIntFromBigInt(task.Payment),
			RequestData:               task.RequestData,
			CallbackAddress:           task.CallbackAddress.Hex(),
			CallbackFunctionId:        task.CallbackFunctionId[:],
			TaskCreatedBlock:          task.TaskCreatedBlock,
			QuorumNumbers:             task.GroupNumbers,
			QuorumThresholdPercentage: task.GroupThresholdPercentage,
			AdvanceDecode:             task.AdvanceDecode,
		},
		PriceFeed: &pricetypes.PriceFeedParam{
			BaseSymbol:  priceFeed.BaseSymbol,
			QuoteSymbol: priceFeed.QuoteSymbol,
		},
	}
	return taskRequest, nil
}

func (d *Dispatcher) serializeScriptTask(chainID uint64, newTask *contractdataoracle.ContractDataOracleNewTaskCreated) (sdk.Msg, error) {
	task := newTask.Task
	scriptData, err := ParseScript(task.RequestData)
	if err != nil {
		d.logger.Error("Failed to parse script data", "chainID", chainID, "error", err)
		return nil, err
	}
	taskRequest := &processortypes.RequestScriptIn{
		TaskIndex:                 newTask.TaskIndex,
		RequestId:                 task.RequestId[:],
		FeeToken:                  task.FeeToken.Hex(),
		Payment:                   math.NewIntFromBigInt(task.Payment),
		RequestData:               task.RequestData,
		CallbackAddress:           task.CallbackAddress.Hex(),
		CallbackFunctionId:        task.CallbackFunctionId[:],
		TaskCreatedBlock:          task.TaskCreatedBlock,
		QuorumNumbers:             task.GroupNumbers,
		QuorumThresholdPercentage: task.GroupThresholdPercentage,
		AdvanceDecode:             task.AdvanceDecode,
		ScriptId:                  scriptData.ScriptId,
		ScriptParam:               scriptData.Params,
	}
	return taskRequest, nil
}

func (d *Dispatcher) serializeVRFTask(chainID uint64, newTask *contractdataoracle.ContractDataOracleNewTaskCreated) (sdk.Msg, error) {
	task := newTask.Task
	param, err := ParseVRFTaskParam(task.RequestData)
	if err != nil {
		d.logger.Error("Failed to parse VRF task param", "chainID", chainID, "error", err)
		return nil, err
	}

	// Gnerate random number
	privKeyStr := d.config.ECCKeyPair.ECCPrivateKey
	if strings.HasPrefix(privKeyStr, "0x") {
		privKeyStr = privKeyStr[2:]
	}
	priKeyBuf, err := hex.DecodeString(privKeyStr)
	if err != nil {
		d.logger.Error("Failed to decode private key", "error", err)
	}
	priKey, err := keypair.DeserializePrivateKey(priKeyBuf)
	if err != nil {
		d.logger.Error("Failed to deserialize private key", "error", err)
	}

	// Generate a VRF
	var vrfData []*types.VRFData

	taskMetadata := &types.TaskMetadata{
		TaskIndex:                newTask.TaskIndex,
		RequestId:                task.RequestId[:],
		FeeToken:                 task.FeeToken.Hex(),
		Payment:                  math.NewIntFromBigInt(task.Payment),
		RequestData:              task.RequestData,
		CallbackAddress:          task.CallbackAddress.Hex(),
		CallbackFunctionId:       task.CallbackFunctionId[:],
		TaskCreatedBlock:         task.TaskCreatedBlock,
		GroupNumbers:             task.GroupNumbers,
		GroupThresholdPercentage: task.GroupThresholdPercentage,
		AdvanceDecode:            task.AdvanceDecode,
	}

	for i := 0; i < int(param.NumWords); i++ {
		vrfValue, vrfProof, err := common.ComputeVRF(priKey, taskMetadata)
		if err != nil {
			d.logger.Error("Failed to generate VRF", "error", err)
			return nil, err
		}
		vrfData = append(vrfData, &types.VRFData{
			VrfValue: vrfValue,
			VrfProof: vrfProof,
		})
	}

	taskRequest := &types.VRFTaskRequest{
		TaskMetadata: taskMetadata,
		VrfData:      vrfData,
	}

	return taskRequest, nil
}

// listenForNewTasks listens for new price tasks on a specific chain
func (d *Dispatcher) dispatchTask(
	taskData []byte,
	groupNumbersBytes []byte,
	chainID uint64,
	blockNumber uint32,
	groupThresholdPercentage uint32,
) error {
	// Log initial parameters
	d.logger.Info("dispatchTask called",
		"taskData", taskData,
		"groupNumbersBytes", groupNumbersBytes,
		"chainID", chainID,
		"blockNumber", blockNumber,
		"groupThresholdPercentage", groupThresholdPercentage,
	)

	// Convert byte slice to group numbers
	groupNumbers := make([]uint32, len(groupNumbersBytes))
	groupNumbersForInteractor := make([]interactortypes.GroupNumber, len(groupNumbersBytes))
	for i, b := range groupNumbersBytes {
		groupNumbers[i] = uint32(b)
		groupNumbersForInteractor[i] = interactortypes.GroupNumber(b)
	}

	// Retrieve operator DVS state
	operatorDVSState, err := d.reader.GetOperatorsDVSStateAtBlock(
		chainID,
		groupNumbersForInteractor,
		blockNumber,
	)
	if err != nil {
		return fmt.Errorf("failed to retrieve operator DVS state: %w", err)
	}

	if len(operatorDVSState) == 0 {
		return errors.New("no operator DVS state found")
	}

	d.logger.Info("operator DVS state retrieved",
		"count", len(operatorDVSState),
		"chainID", chainID,
	)

	// Iterate over each operator's DVS state
	for operatorID, operatorState := range operatorDVSState {
		// Fetch operator info
		info, err := d.reader.GetOperatorInfoByID(operatorID)
		if err != nil {
			d.logger.Error("failed to get operator info",
				"chainID", chainID,
				"error", err,
				"operatorID", operatorID,
				"operatorAddress", operatorState.OperatorAddress,
			)
			continue
		}

		d.logger.Info("preparing to send task to DVS operator",
			"chainID", chainID,
			"operatorID", operatorID,
			"operatorAddress", operatorState.OperatorAddress,
			"socket", info.Socket,
		)

		// Create the client
		client, err := http.New(info.Socket.String(), "")
		if err != nil {
			d.logger.Error("failed to create eth client",
				"chainID", chainID,
				"error", err,
				"operatorID", operatorID,
				"operatorAddress", operatorState.OperatorAddress,
				"socket", info.Socket,
			)
			continue
		}

		// Send the task asynchronously
		reqResp, err := client.RequestDVSAsync(
			context.Background(),
			taskData,
			int64(blockNumber),
			int64(chainID),
			groupNumbers,
			[]uint32{groupThresholdPercentage},
		)
		if err != nil {
			d.logger.Error("failed to send task to PellDVS",
				"chainID", chainID,
				"error", err,
				"operatorID", operatorID,
				"operatorAddress", operatorState.OperatorAddress,
				"socket", info.Socket,
			)
			continue
		}

		d.logger.Info("task sent to PellDVS successfully",
			"chainID", chainID,
			"operatorID", operatorID,
			"operatorAddress", operatorState.OperatorAddress,
			"socket", info.Socket,
			"response", reqResp,
		)
	}

	return nil
}
