package dispatcher

import (
	"cosmossdk.io/math"
	"fmt"
	contractdataoracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"intellix/dvs"
	pricetypes "intellix/dvs/price/types"
	processortypes "intellix/dvs/processor/types"
)

func (d *Dispatcher) listenForNewPriceTasks(chain *chainWatcher) {
	d.logger.Info("now listen for new tasks")
	newTaskChan := make(chan *contractdataoracle.ContractDataOracleNewTaskCreated)
	// TODO: add index by height
	// TODO: add scan mode
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
			d.handleNewPriceTask(chain.chainID, newTask)
		case <-d.Quit():
			return
		}
	}
}

func (d *Dispatcher) handleNewPriceTask(chainID uint64, newTask *contractdataoracle.ContractDataOracleNewTaskCreated) {
	d.logger.Info("New task created",
		"chainID", chainID,
		"TaskIndex", newTask.TaskIndex,
		"RequestId", newTask.Task.RequestId,
		"TaskType", newTask.Task.TaskType,
	)

	taskData, err := d.serializeNewPriceTask(chainID, newTask)
	if err != nil {
		d.logger.Error("Failed to serialize task", "chainID", chainID, "error", err)
		return
	}

	err = d.dispatchTask(taskData, newTask.Task.GroupNumbers, chainID, uint32(newTask.Raw.BlockNumber), newTask.Task.GroupThresholdPercentage)
	if err != nil {
		d.logger.Error("Failed to dispatch task", "chainID", chainID, "error", err)
	}
}

func (d *Dispatcher) serializeNewPriceTask(chainID uint64, newTask *contractdataoracle.ContractDataOracleNewTaskCreated) ([]byte, error) {
	d.logger.Info("serializeNewPriceTask",
		"chainID", chainID,
		"taskIndex", newTask.TaskIndex,
		"task", fmt.Sprintf("%+v", newTask.Task),
	)

	task := newTask.Task
	var taskRequest sdk.Msg

	if task.TaskType.Int64() == dvs.TaskTypePriceFeed {
		priceFeed, err := ParsePriceFeed(newTask.Task.RequestData)
		if err != nil {
			d.logger.Error("Failed to parse price feed", "chainID", chainID, "error", err)
			return nil, err
		}
		taskRequest = &pricetypes.RequestPriceFeedIn{
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
	} else if task.TaskType.Int64() == dvs.TaskTypeProcessor {
		scriptData, err := ParseScript(newTask.Task.RequestData)
		if err != nil {
			d.logger.Error("Failed to parse script data", "chainID", chainID, "error", err)
			return nil, err
		}
		taskRequest = &processortypes.RequestScriptIn{
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
	}
	if taskRequest == nil {
		return nil, fmt.Errorf("invalid task request")
	}

	return d.msgEncoder.EncodeMsgs(taskRequest)
}
