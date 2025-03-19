package taskdispatcher

import "C"

import (
	"context"
	"fmt"
	"sync"

	"cosmossdk.io/math"
	"github.com/0xPellNetwork/pellapp-sdk/service/tx"
	interactorcfg "github.com/0xPellNetwork/pelldvs-interactor/config"
	"github.com/0xPellNetwork/pelldvs-interactor/interactor/reader"
	interactortypes "github.com/0xPellNetwork/pelldvs-interactor/types"
	dvslog "github.com/0xPellNetwork/pelldvs-libs/log"
	pelldvscfg "github.com/0xPellNetwork/pelldvs/config"
	"github.com/0xPellNetwork/pelldvs/rpc/client/http"
	rpclocal "github.com/0xPellNetwork/pelldvs/rpc/client/local"
	contractdataoracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/libs/service"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"intellix/config"
	pricetypes "intellix/dvs/price/types"
	processortypes "intellix/dvs/processor/types"
)

type TaskDispatcher struct {
	service.BaseService

	logger        dvslog.Logger
	pellDVSClient *rpclocal.Local
	chains        map[uint64]*chainWatcher
	mu            sync.Mutex
	msgEncoder    tx.MsgEncoder
	reader        *reader.DVSReaderServer
}

type chainWatcher struct {
	chainID  uint64
	contract *contractdataoracle.ContractDataOracle
	client   *ethclient.Client
}

func newTaskProtoEncoder() tx.MsgEncoder {
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	return tx.NewDefaultDecoder(cdc)
}

func NewTaskDispatcher(logger dvslog.Logger, pellDVSConf *pelldvscfg.Config, tdConf *Config) (*TaskDispatcher, error) {
	logger = logger.With("comp", "dispatcher")
	// load interactor config
	iteractorConfig, err := interactorcfg.LoadConfig(pellDVSConf.Pell.InteractorConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load interactor configuration")
	}

	//  create db
	db, err := pelldvscfg.DefaultDBProvider(&pelldvscfg.DBContext{
		ID:     "indexer",
		Config: pellDVSConf,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %v", err)
	}

	td := &TaskDispatcher{
		logger:     logger,
		chains:     make(map[uint64]*chainWatcher),
		msgEncoder: newTaskProtoEncoder(),
	}

	for _, chainConfig := range tdConf.Chains {
		if err := td.AddChain(chainConfig); err != nil {
			return nil, fmt.Errorf("failed to add chain %d: %w", chainConfig.ChainID, err)
		}
	}

	dvsReader, err := reader.NewDVSReaderFromConfig(iteractorConfig, db, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create DVS reader: %w", err)
	}
	td.reader = dvsReader
	td.BaseService = *service.NewBaseService(nil, "TaskDispatcher", td)
	return td, nil
}

func (td *TaskDispatcher) AddChain(config config.ChainConfig) error {
	td.mu.Lock()
	defer td.mu.Unlock()
	td.logger.Info(fmt.Sprintf("listen chain, chainID: %d, url: %s, address: %s", config.ChainID, config.RPCURL, config.ContractAddress))

	if err := config.Validate(); err != nil {
		return err
	}
	if _, exists := td.chains[config.ChainID]; exists {
		return fmt.Errorf("chain %d already exists", config.ChainID)
	}

	wsClient, err := ethclient.Dial(config.WSURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	contract, err := contractdataoracle.NewContractDataOracle(common.HexToAddress(config.ContractAddress), wsClient)
	if err != nil {
		return fmt.Errorf("failed to instantiate contract: %w", err)
	}

	td.chains[config.ChainID] = &chainWatcher{
		chainID:  config.ChainID,
		contract: contract,
		client:   wsClient,
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
	newTaskChan := make(chan *contractdataoracle.ContractDataOracleNewTaskCreated)
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

func (td *TaskDispatcher) handleNewTask(chainID uint64, newTask *contractdataoracle.ContractDataOracleNewTaskCreated) {
	td.logger.Info("New task created",
		"chainID", chainID,
		"TaskIndex", newTask.TaskIndex,
		"RequestId", newTask.Task.RequestId,
		"TaskType", newTask.Task.TaskType,
	)

	taskData, err := td.serializeTask(chainID, newTask)
	if err != nil {
		td.logger.Error("Failed to serialize task", "chainID", chainID, "error", err)
		return
	}

	quorumNumbers := make([]uint32, len(newTask.Task.GroupNumbers))
	quorumNumbersForInteractor := make([]interactortypes.GroupNumber, len(newTask.Task.GroupNumbers))
	for i, b := range newTask.Task.GroupNumbers {
		quorumNumbers[i] = uint32(b)
		quorumNumbersForInteractor[i] = interactortypes.GroupNumber(b)
	}

	operatorDVSState, err := td.reader.GetOperatorsDVSStateAtBlock(chainID,
		quorumNumbersForInteractor,
		uint32(newTask.Raw.BlockNumber),
	)
	if err != nil {
		td.logger.Error("Failed to get operator DVS state",
			"chainID", chainID,
			"blockNumber", newTask.Raw.BlockNumber,
			"quorumNumbers", quorumNumbers,
			"error", err,
		)
		return
	}

	if len(operatorDVSState) == 0 {
		td.logger.Error("No operator DVS state found",
			"chainID", chainID,
			"blockNumber", newTask.Raw.BlockNumber,
			"quorumNumbers", quorumNumbers,
			"error", err,
		)
		return
	}

	td.logger.Info("Operator DVS state count", "count", len(operatorDVSState))

	for operatorID, operatorState := range operatorDVSState {
		info, err := td.reader.GetOperatorInfoByID(operatorID)
		if err != nil {
			td.logger.Error("Failed to get operator info",
				"chainID", chainID,
				"error", err,
				"operatorID", operatorID,
				"operatorAddress", operatorState.OperatorAddress,
			)
			continue
		}

		td.logger.Info("prepare to send task to DVS app operator",
			"chainID", chainID,
			"TaskIndex", newTask.TaskIndex,
			"RequestId", newTask.Task.RequestId,
			"operatorID", operatorID,
			"operatorAddress", operatorState.OperatorAddress,
			"socket", info.Socket,
		)

		client, err := http.New(info.Socket.String(), "")
		if err != nil {
			td.logger.Error("Failed to create eth client",
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
			quorumNumbers,
			[]uint32{newTask.Task.GroupThresholdPercentage},
		)
		if err != nil {
			td.logger.Error("Failed to send task to PellDVS",
				"chainID", chainID, "error", err,
				"operatorID", operatorID,
				"operatorAddress", operatorState.OperatorAddress,
				"socket", info.Socket,
			)
			continue
		}

		td.logger.Info("Task sent to PellDVS successfully",
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

func (td *TaskDispatcher) serializeTask(chainID uint64, newTask *contractdataoracle.ContractDataOracleNewTaskCreated) ([]byte, error) {
	td.logger.Info("serializeTask",
		"chainID", chainID,
		"taskIndex", newTask.TaskIndex,
		"task", fmt.Sprintf("%+v", newTask.Task),
	)

	task := newTask.Task
	var taskRequest sdk.Msg

	if task.TaskType.Int64() == TaskTypePrice {
		priceFeed, err := ParsePriceFeed(newTask.Task.RequestData)
		if err != nil {
			td.logger.Error("Failed to parse price feed", "chainID", chainID, "error", err)
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
	} else if task.TaskType.Int64() == TaskTypeScript {
		scriptData, err := ParseScript(newTask.Task.RequestData)
		if err != nil {
			td.logger.Error("Failed to parse script data", "chainID", chainID, "error", err)
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
