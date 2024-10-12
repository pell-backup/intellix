package taskdispatcher

import (
	"context"
	"fmt"
	avsi "github.com/0xPellNetwork/pelldvs/application"
	"github.com/0xPellNetwork/pelldvs/avsi/types"
	contractPriceOracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/PriceOracle"
	"github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/libs/service"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"intellix/pkg/pelldvs"
	"sync"
)

type TaskDispatcher struct {
	service.BaseService

	logger        log.Logger
	pellDVSClient *pelldvs.Client
	chains        map[string]*chainWatcher
	mu            sync.Mutex
}

type chainWatcher struct {
	chainID  string
	contract *contractPriceOracle.ContractPriceOracle
	client   *ethclient.Client
}

func NewTaskDispatcher(logger log.Logger, pellDVSClient *pelldvs.Client, configs []*ChainConfig) (*TaskDispatcher, error) {
	td := &TaskDispatcher{
		logger:        logger,
		pellDVSClient: pellDVSClient,
		chains:        make(map[string]*chainWatcher),
	}

	for _, config := range configs {
		if err := td.AddChain(config); err != nil {
			return nil, fmt.Errorf("failed to add chain %s: %w", config.ChainID, err)
		}
	}

	td.BaseService = *service.NewBaseService(logger, "TaskDispatcher", td)
	return td, nil
}

func (td *TaskDispatcher) AddChain(config *ChainConfig) error {
	td.mu.Lock()
	defer td.mu.Unlock()

	if err := config.Validate(); err != nil {
		return err
	}
	if _, exists := td.chains[config.ChainID]; exists {
		return fmt.Errorf("chain %s already exists", config.ChainID)
	}

	ethClient, err := ethclient.Dial(config.EthURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	contract, err := contractPriceOracle.NewContractPriceOracle(common.HexToAddress(config.ContractAddress), ethClient)
	if err != nil {
		return fmt.Errorf("failed to instantiate contract: %w", err)
	}

	td.chains[config.ChainID] = &chainWatcher{
		chainID:  config.ChainID,
		contract: contract,
		client:   ethClient,
	}

	return nil
}

func (td *TaskDispatcher) Start() error {
	for _, chain := range td.chains {
		go td.listenForNewTasks(chain)
	}
	return nil
}

func (td *TaskDispatcher) listenForNewTasks(chain *chainWatcher) {
	newTaskChan := make(chan *contractPriceOracle.ContractPriceOracleNewTaskCreated)
	sub, err := chain.contract.WatchNewTaskCreated(&bind.WatchOpts{}, newTaskChan, nil)
	if err != nil {
		td.logger.Error("Failed to watch for new tasks", "chainID", chain.chainID, "error", err)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			td.logger.Error("Subscription error", "chainID", chain.chainID, "error", err)
			return
		case newTask := <-newTaskChan:
			td.handleNewTask(chain.chainID, newTask)
		case <-td.Quit():
			return
		}
	}
}

func (td *TaskDispatcher) handleNewTask(chainID string, newTask *contractPriceOracle.ContractPriceOracleNewTaskCreated) {
	td.logger.Info("New task created", "chainID", chainID, "TaskIndex", newTask.TaskIndex, "RequestId", newTask.Task.RequestId)

	taskData, err := td.serializeTask(newTask.Task)
	if err != nil {
		td.logger.Error("Failed to serialize task", "chainID", chainID, "error", err)
		return
	}

	err = td.pellDVSClient.RequestDVS(context.Background(), &avsi.RequestProcessRequest{
		Request: types.DVSRequest{
			Data:    taskData,
			Height:  0,
			ChainID: common.HexToHash(chainID).Big(),
		},
	})
	if err != nil {
		td.logger.Error("Failed to send task to PellDVS", "chainID", chainID, "error", err)
		return
	}

	td.logger.Info("Task sent to PellDVS successfully", "chainID", chainID, "TaskIndex", newTask.TaskIndex)
}

func (td *TaskDispatcher) serializeTask(task contractPriceOracle.IPriceOracleTask) ([]byte, error) {
	// TODO: serialize to proto-buffer, mock json for now
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
