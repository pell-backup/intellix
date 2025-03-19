package dispatcher

import (
	"context"
	"fmt"
	interactortypes "github.com/0xPellNetwork/pelldvs-interactor/types"
	"github.com/0xPellNetwork/pelldvs/rpc/client/http"
	contractdataoracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"intellix/dvs/vrf/types"
)

func (d *Dispatcher) listenForNewVRFTasks(chain *chainWatcher) {
	d.logger.Info("now listen for new tasks")
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
			d.handleNewVRFTask(chain.chainID, newTask)
		case <-d.Quit():
			return
		}
	}
}

func (d *Dispatcher) handleNewVRFTask(chainID uint64, newTask *contractdataoracle.ContractDataOracleNewTaskCreated) {
	d.logger.Info("New task created",
		"chainID", chainID,
		"TaskIndex", newTask.TaskIndex,
		"RequestId", newTask.Task.RequestId,
		"TaskType", newTask.Task.TaskType,
	)

	taskData, err := d.serializeNewVRFTask(chainID, newTask)
	if err != nil {
		d.logger.Error("Failed to serialize task", "chainID", chainID, "error", err)
		return
	}

	groupNumbers := make([]uint32, len(newTask.Task.GroupNumbers))
	groupNumbersForInteractor := make([]interactortypes.GroupNumber, len(newTask.Task.GroupNumbers))
	for i, b := range newTask.Task.GroupNumbers {
		groupNumbers[i] = uint32(b)
		groupNumbersForInteractor[i] = interactortypes.GroupNumber(b)
	}

	operatorDVSState, err := d.reader.GetOperatorsDVSStateAtBlock(chainID,
		groupNumbersForInteractor,
		uint32(newTask.Raw.BlockNumber),
	)
	if err != nil {
		d.logger.Error("Failed to get operator DVS state",
			"chainID", chainID,
			"blockNumber", newTask.Raw.BlockNumber,
			"groupNumbers", groupNumbers,
			"error", err,
		)
		return
	}

	if len(operatorDVSState) == 0 {
		d.logger.Error("No operator DVS state found",
			"chainID", chainID,
			"blockNumber", newTask.Raw.BlockNumber,
			"groupNumbers", groupNumbers,
			"error", err,
		)
		return
	}

	d.logger.Info("Operator DVS state count", "count", len(operatorDVSState))

	for operatorID, operatorState := range operatorDVSState {
		info, err := d.reader.GetOperatorInfoByID(operatorID)
		if err != nil {
			d.logger.Error("Failed to get operator info",
				"chainID", chainID,
				"error", err,
				"operatorID", operatorID,
				"operatorAddress", operatorState.OperatorAddress,
			)
			continue
		}

		d.logger.Info("prepare to send task to DVS app operator",
			"chainID", chainID,
			"TaskIndex", newTask.TaskIndex,
			"RequestId", newTask.Task.RequestId,
			"operatorID", operatorID,
			"operatorAddress", operatorState.OperatorAddress,
			"socket", info.Socket,
		)

		client, err := http.New(info.Socket.String(), "")
		if err != nil {
			d.logger.Error("Failed to create eth client",
				"chainID", chainID,
				"error", err,
				"operatorID", operatorID,
				"operatorAddress", operatorState.OperatorAddress,
				"socket", info.Socket,
			)
			continue
		}
		reqResp, err := client.RequestDVSAsync(
			context.Background(),
			taskData,
			int64(newTask.Raw.BlockNumber),
			int64(chainID),
			groupNumbers,
			[]uint32{newTask.Task.GroupThresholdPercentage},
		)
		if err != nil {
			d.logger.Error("Failed to send task to PellDVS",
				"chainID", chainID, "error", err,
				"operatorID", operatorID,
				"operatorAddress", operatorState.OperatorAddress,
				"socket", info.Socket,
			)
			continue
		}

		d.logger.Info("Task sent to PellDVS successfully",
			"chainID", chainID,
			"TaskIndex", newTask.TaskIndex,
			"RequestId", newTask.Task.RequestId,
			"operatorID", operatorID,
			"operatorAddress", operatorState.OperatorAddress,
			"socket", info.Socket,
			"response", reqResp,
		)

	}
}

func (d *Dispatcher) serializeNewVRFTask(chainID uint64, task *contractdataoracle.ContractDataOracleNewTaskCreated) ([]byte, error) {
	d.logger.Info("serializeNewVRFTask", "chainID", chainID, "taskIndex", task.TaskIndex, "task", fmt.Sprintf("%+v", task.Task))
	taskRequest := &types.GenerateRandomNumberRequest{
		&types.Task{
			TaskIndex:                task.TaskIndex,
			Height:                   task.Task.TaskCreatedBlock,
			ChainId:                  chainID,
			GroupNumbers:             task.Task.GroupNumbers,
			GroupThresholdPercentage: task.Task.GroupThresholdPercentage,
		},
	}
	return d.msgEncoder.EncodeMsgs(taskRequest)
}
