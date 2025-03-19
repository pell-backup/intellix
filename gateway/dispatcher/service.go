package dispatcher

import "C"

import (
	"context"
	"fmt"
	interactortypes "github.com/0xPellNetwork/pelldvs-interactor/types"
	"github.com/0xPellNetwork/pelldvs/rpc/client/http"
	"intellix/gateway/types"
	"sync"

	interactorcfg "github.com/0xPellNetwork/pelldvs-interactor/config"
	"github.com/0xPellNetwork/pelldvs-interactor/interactor/reader"
	dvslog "github.com/0xPellNetwork/pelldvs-libs/log"
	pelldvscfg "github.com/0xPellNetwork/pelldvs/config"
	rpclocal "github.com/0xPellNetwork/pelldvs/rpc/client/local"
	contractdataoracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/libs/service"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"intellix/sdk/dvs_msg_handler/tx"
)

type Dispatcher struct {
	service.BaseService

	logger        dvslog.Logger
	pellDVSClient *rpclocal.Local
	chains        map[uint64]*chainWatcher
	mu            sync.Mutex
	msgEncoder    tx.MsgEncoder
	reader        *reader.DVSReaderServer
}

type chainWatcher struct {
	chainID            uint64
	dataOracleContract *contractdataoracle.ContractDataOracle
	client             *ethclient.Client
}

func newTaskProtoEncoder() tx.MsgEncoder {
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	return tx.NewDefaultDecoder(cdc)
}

func NewDispatcher(logger dvslog.Logger, pellDVSConf *pelldvscfg.Config, conf *types.Config) (*Dispatcher, error) {
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

	dispatcher := &Dispatcher{
		logger:     logger,
		chains:     make(map[uint64]*chainWatcher),
		msgEncoder: newTaskProtoEncoder(),
	}

	for _, chainConfig := range conf.Chains {
		if err := dispatcher.AddChain(chainConfig); err != nil {
			return nil, fmt.Errorf("failed to add chain %d: %w", chainConfig.ChainID, err)
		}
	}

	dvsReader, err := reader.NewDVSReaderFromConfig(iteractorConfig, db, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create DVS reader: %w", err)
	}

	dispatcher.reader = dvsReader
	dispatcher.BaseService = *service.NewBaseService(nil, "Dispatcher", dispatcher)
	return dispatcher, nil
}

func (d *Dispatcher) AddChain(config types.ChainConfig) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.logger.Info(fmt.Sprintf("listen chain, chainID: %d, url: %s, address: %s", config.ChainID, config.RPCURL, config.ContractAddress))

	if err := config.Validate(); err != nil {
		return err
	}
	if _, exists := d.chains[config.ChainID]; exists {
		return fmt.Errorf("chain %d already exists", config.ChainID)
	}

	wsClient, err := ethclient.Dial(config.WSURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	contract, err := contractdataoracle.NewContractDataOracle(common.HexToAddress(config.ContractAddress), wsClient)
	if err != nil {
		return fmt.Errorf("failed to instantiate dataOracleContract: %w", err)
	}

	d.chains[config.ChainID] = &chainWatcher{
		chainID:            config.ChainID,
		dataOracleContract: contract,
		client:             wsClient,
	}

	return nil
}

// listenForNewPriceTasks listens for new price tasks on a specific chain
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

func (d *Dispatcher) Start() error {
	for _, chain := range d.chains {
		go d.listenForNewPriceTasks(chain)
	}
	return nil
}

func (d *Dispatcher) OnStart() error {
	return nil
}

func (d *Dispatcher) Stop() error {
	return nil
}

func (d *Dispatcher) OnStop() {
}

func (d *Dispatcher) Reset() error {
	return nil
}

func (d *Dispatcher) OnReset() error {
	return nil
}

func (d *Dispatcher) IsRunning() bool {
	return d.BaseService.IsRunning()
}

func (d *Dispatcher) Quit() <-chan struct{} {
	return d.BaseService.Quit()
}

func (d *Dispatcher) String() string {
	return d.BaseService.String()
}

func (d *Dispatcher) SetLogger(logger log.Logger) {
	d.BaseService.SetLogger(logger)
}
