package taskdispatcher

import "C"
import (
	"context"
	"fmt"
	"intellix/pkg/dvs_msg_handler/tx"
	"intellix/pkg/pelldvs"
	pricetypes "intellix/x/price/dvs/types"
	"sync"

	dvslog "github.com/0xPellNetwork/pelldvs/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/math"
	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
	contractDataOracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/libs/service"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TaskDispatcher struct {
	service.BaseService

	logger        dvslog.Logger
	pellDVSClient *pelldvs.Client
	chains        map[uint64]*chainWatcher
	mu            sync.Mutex
	msgEncoder    tx.MsgEncoder
}

type chainWatcher struct {
	chainID  uint64
	contract *contractDataOracle.ContractDataOracle
	client   *ethclient.Client
}

func newTaskProtoEncoder() tx.MsgEncoder {
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	return tx.NewDefaultDecoder(cdc)
}

func NewTaskDispatcher(logger dvslog.Logger, pellDVSClient *pelldvs.Client, configs []*ChainConfig) (*TaskDispatcher, error) {
	td := &TaskDispatcher{
		logger:        logger,
		pellDVSClient: pellDVSClient,
		chains:        make(map[uint64]*chainWatcher),
		msgEncoder:    newTaskProtoEncoder(),
	}

	for _, config := range configs {
		if err := td.AddChain(config); err != nil {
			return nil, fmt.Errorf("failed to add chain %d: %w", config.ChainID, err)
		}
	}

	td.BaseService = *service.NewBaseService(nil, "TaskDispatcher", td)
	return td, nil
}

func (td *TaskDispatcher) AddChain(config *ChainConfig) error {
	td.mu.Lock()
	defer td.mu.Unlock()
	td.logger.Info(fmt.Sprintf("listen chain, chainID: %d, url: %s, address: %s", config.ChainID, config.EthURL, config.ContractAddress))

	if err := config.Validate(); err != nil {
		return err
	}
	if _, exists := td.chains[config.ChainID]; exists {
		return fmt.Errorf("chain %d already exists", config.ChainID)
	}

	ethClient, err := ethclient.Dial(config.EthURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	contract, err := contractDataOracle.NewContractDataOracle(common.HexToAddress(config.ContractAddress), ethClient)
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
	td.logger.Info("now listen for new tasks")
	newTaskChan := make(chan *contractDataOracle.ContractDataOracleNewTaskCreated)
	// TODO: add index by height
	// TODO: add scan mode
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

func (td *TaskDispatcher) handleNewTask(chainID uint64, newTask *contractDataOracle.ContractDataOracleNewTaskCreated) {
	td.logger.Info("New task created", "chainID", chainID, "TaskIndex", newTask.TaskIndex, "RequestId", newTask.Task.RequestId, "TaskType", newTask.Task.TaskType)

	taskData, err := td.serializeTask(chainID, newTask)
	if err != nil {
		td.logger.Error("Failed to serialize task", "chainID", chainID, "error", err)
		return
	}

	quorumNumbers := make([]uint32, len(newTask.Task.GroupNumbers))
	for i, b := range newTask.Task.GroupNumbers {
		quorumNumbers[i] = uint32(b)
	}

	err = td.pellDVSClient.RequestDVS(context.Background(), &avsitypes.RequestProcessDVSRequest{
		Request: &avsitypes.DVSRequest{
			Data:                      taskData,
			Height:                    int64(newTask.Raw.BlockNumber),
			ChainId:                   int64(chainID),
			GroupNumbers:              quorumNumbers,
			GroupThresholdPercentages: []uint32{newTask.Task.GroupThresholdPercentage},
		},
	})
	if err != nil {
		td.logger.Error("Failed to send task to PellDVS", "chainID", chainID, "error", err)
		return
	}

	td.logger.Info("Task sent to PellDVS successfully", "chainID", chainID, "TaskIndex", newTask.TaskIndex)
}

func (td *TaskDispatcher) serializeTask(chainID uint64, newTask *contractDataOracle.ContractDataOracleNewTaskCreated) ([]byte, error) {
	td.logger.Info("serializeTask",
		"chainID", chainID,
		"taskIndex", newTask.TaskIndex,
		"task", fmt.Sprintf("%+v", newTask.Task),
	)
	priceFeed, err := ParsePriceFeed(newTask.Task.RequestData)
	if err != nil {
		td.logger.Error("Failed to parse price feed", "chainID", chainID, "error", err)
		return nil, err
	}
	task := newTask.Task
	var taskRequest sdk.Msg

	// TODO: add more task-types
	if task.TaskType.Int64() == TaskTypePrice {
		taskRequest = &pricetypes.ProcessRequestPriceFeedIn{
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
			},
			PriceFeed: &pricetypes.PriceFeedParam{
				BaseSymbol:  priceFeed.BaseSymbol,
				QuoteSymbol: priceFeed.QuoteSymbol,
			},
		}
	}
	if taskRequest == nil {
		return nil, fmt.Errorf("invalid task request")
	}

	return td.msgEncoder.EncodeMsgs(taskRequest)
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
