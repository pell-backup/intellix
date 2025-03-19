package dispatcher

import (
	"context"
	"encoding/hex"
	"fmt"
	interactortypes "github.com/0xPellNetwork/pelldvs-interactor/types"
	"github.com/0xPellNetwork/pelldvs/rpc/client/http"
	contractdataoracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ontio/ontology-crypto/keypair"
	"intellix/dvs/vrf/types"
	"strings"
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

func (d *Dispatcher) handleNewVRFTask(chainID uint64, task *contractdataoracle.ContractDataOracleNewTaskCreated) {
	d.logger.Info("New task created",
		"chainID", chainID,
		"TaskIndex", task.TaskIndex,
		"RequestId", task.Task.RequestId,
		"TaskType", task.Task.TaskType,
	)

	// Get private key
	// TODO: use config
	privKeyStr := "0x1506060172d45fe6b916c5f903be549d11543140bb7b28a306483e09bf14dce26c"
	// pubkey := 1504e7cae432489d7c670139a3cb712f0766c66e52525242e5d7672f11b6c8c1f6bd5f859bd281955549ee5dc697a12028189fac013d9a7b1c4c9aab00823d3d18e6
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
	vrfValue, vrfProof, err := computeVrf(priKey, task)
	if err != nil {
		d.logger.Error("Failed to generate VRF", "error", err)
		return
	}

	// Serialize the task with the VRF
	taskData, err := d.serializeNewVRFTask(chainID, task, vrfValue, vrfProof)
	if err != nil {
		d.logger.Error("Failed to serialize task", "chainID", chainID, "error", err)
		return
	}

	groupNumbers := make([]uint32, len(task.Task.GroupNumbers))
	groupNumbersForInteractor := make([]interactortypes.GroupNumber, len(task.Task.GroupNumbers))
	for i, b := range task.Task.GroupNumbers {
		groupNumbers[i] = uint32(b)
		groupNumbersForInteractor[i] = interactortypes.GroupNumber(b)
	}

	operatorDVSState, err := d.reader.GetOperatorsDVSStateAtBlock(
		chainID,
		groupNumbersForInteractor,
		uint32(task.Raw.BlockNumber),
	)
	if err != nil {
		d.logger.Error("Failed to get operator DVS state",
			"chainID", chainID,
			"blockNumber", task.Raw.BlockNumber,
			"groupNumbers", groupNumbers,
			"error", err,
		)
		return
	}

	if len(operatorDVSState) == 0 {
		d.logger.Error("No operator DVS state found",
			"chainID", chainID,
			"blockNumber", task.Raw.BlockNumber,
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
			"TaskIndex", task.TaskIndex,
			"RequestId", task.Task.RequestId,
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
			int64(task.Raw.BlockNumber),
			int64(chainID),
			groupNumbers,
			[]uint32{task.Task.GroupThresholdPercentage},
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
			"TaskIndex", task.TaskIndex,
			"RequestId", task.Task.RequestId,
			"operatorID", operatorID,
			"operatorAddress", operatorState.OperatorAddress,
			"socket", info.Socket,
			"response", reqResp,
		)

	}
}

func (d *Dispatcher) serializeNewVRFTask(chainID uint64, task *contractdataoracle.ContractDataOracleNewTaskCreated,
	vrfValue, vrfProof []byte) ([]byte, error) {

	d.logger.Info("serializeNewVRFTask", "chainID", chainID, "taskIndex", task.TaskIndex,
		"task", fmt.Sprintf("%+v", task.Task), "vrfValue", hex.EncodeToString(vrfValue), "vrfProof", hex.EncodeToString(vrfProof))

	taskRequest := &types.GenerateRandomNumberRequest{
		Task: &types.Task{
			TaskIndex:                task.TaskIndex,
			Height:                   task.Task.TaskCreatedBlock,
			ChainId:                  chainID,
			GroupNumbers:             task.Task.GroupNumbers,
			GroupThresholdPercentage: task.Task.GroupThresholdPercentage,
		},
		VrfValue: vrfValue,
		VrfProof: vrfProof,
	}
	return d.msgEncoder.EncodeMsgs(taskRequest)
}
