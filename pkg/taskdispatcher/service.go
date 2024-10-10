package taskdispatcher

import (
	"context"
	"fmt"
	contractPriceOracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/PriceOracle"
	"github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/libs/service"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"intellix/pkg/pelldvs"
)

type TaskDispatcher struct {
	service.BaseService

	logger        log.Logger
	pellDVSClient *pelldvs.Client
	contract      *contractPriceOracle.ContractPriceOracle
}

func NewTaskDispatcher(logger log.Logger, ethURL, contractAddress string, pellDVSClient *pelldvs.Client) (*TaskDispatcher, error) {
	ethClient, err := ethclient.Dial(ethURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	contract, err := contractPriceOracle.NewContractPriceOracle(common.HexToAddress(contractAddress), ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate contract: %w", err)
	}

	td := &TaskDispatcher{
		logger:        logger,
		contract:      contract,
		pellDVSClient: pellDVSClient,
	}

	td.BaseService = *service.NewBaseService(logger, "TaskDispatcher", td)
	return td, nil
}

func (td *TaskDispatcher) listenForNewTasks() {
	newTaskChan := make(chan *contractPriceOracle.ContractPriceOracleNewTaskCreated)
	sub, err := td.contract.WatchNewTaskCreated(&bind.WatchOpts{}, newTaskChan, nil)
	if err != nil {
		td.logger.Error("Failed to watch for new tasks", "error", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			td.logger.Error("Subscription error", "error", err)
			return
		case newTask := <-newTaskChan:
			td.handleNewTask(newTask)
		case <-td.Quit():
			return
		}
	}
}

func (td *TaskDispatcher) handleNewTask(newTask *contractPriceOracle.ContractPriceOracleNewTaskCreated) {
	td.logger.Info("New task created", "TaskIndex", newTask.TaskIndex, "RequestId", newTask.Task.RequestId)

	// serialize
	taskData, err := td.serializeTask(newTask.Task)
	if err != nil {
		td.logger.Error("Failed to serialize task", "error", err)
		return
	}

	// pell-dvs client
	err = td.pellDVSClient.RequestDVS(context.Background(), taskData)
	if err != nil {
		td.logger.Error("Failed to send task to PellDVS", "error", err)
		return
	}

	td.logger.Info("Task sent to PellDVS successfully", "TaskIndex", newTask.TaskIndex)
}

func (td *TaskDispatcher) serializeTask(task contractPriceOracle.IPriceOracleTask) ([]byte, error) {
	// todo: serialize to proto-buffer, mock json for now
	return json.Marshal(map[string]interface{}{
		"RequestId":                 task.RequestId,
		"requestData":               task.RequestData,
		"callbackAddress":           task.CallbackAddress.Hex(),
		"callbackFunctionId":        task.CallbackFunctionId,
		"taskCreatedBlock":          task.TaskCreatedBlock,
		"quorumNumbers":             task.QuorumNumbers,
		"quorumThresholdPercentage": task.QuorumThresholdPercentage,
	})
}

func (td *TaskDispatcher) Start() error {
	go td.listenForNewTasks()
	return nil
}

func (td *TaskDispatcher) OnStart() error {
	return nil
}

func (td *TaskDispatcher) Stop() error {
	return nil
}

func (td *TaskDispatcher) OnStop() {
}

func (td *TaskDispatcher) Reset() error {
	return nil
}

func (td *TaskDispatcher) OnReset() error {
	return nil
}

func (td *TaskDispatcher) IsRunning() bool {
	return td.BaseService.IsRunning()
}

func (td *TaskDispatcher) Quit() <-chan struct{} {
	return td.BaseService.Quit()
}

func (td *TaskDispatcher) String() string {
	return td.BaseService.String()
}

func (td *TaskDispatcher) SetLogger(logger log.Logger) {
	td.BaseService.SetLogger(logger)
}
