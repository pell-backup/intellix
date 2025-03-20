package submitter

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync"

	"github.com/0xPellNetwork/pelldvs-libs/log"
	contractdataoracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/cometbft/cometbft/libs/service"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"intellix/gateway/types"
)

type ChainConnection struct {
	ethClient          *ethclient.Client
	contractDataOracle *contractdataoracle.ContractDataOracle
}

type Submitter struct {
	service.BaseService

	server     *rpc.Server
	serverAddr string
	listener   net.Listener

	cfg *types.Config
	ctx context.Context

	logger     log.Logger
	privateKey *keystore.Key

	chainConnections map[uint64]*ChainConnection
	taskMap          sync.Map
	nonceMap         sync.Map
}

func NewSubmitter(logger log.Logger, ctx context.Context, cfg *types.Config) (*Submitter, error) {
	logger = logger.With("comp", "submitter")

	chainConns := make(map[uint64]*ChainConnection)
	for chainID, chainCfg := range cfg.Chains {
		ethClient, err := ethclient.Dial(chainCfg.RPCURL)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to chain %d: %v", chainID, err)
		}

		contract, err := contractdataoracle.NewContractDataOracle(common.HexToAddress(chainCfg.ContractAddress), ethClient)
		if err != nil {
			return nil, fmt.Errorf("failed to create contract instance for chain %d: %v", chainID, err)
		}

		chainConns[chainID] = &ChainConnection{
			ethClient:          ethClient,
			contractDataOracle: contract,
		}
	}

	// Read private key
	keyJSON, err := os.ReadFile(cfg.PrivateKeyStorePath)
	if err != nil {
		logger.Error("Failed to read private key file", "error", err)
		return nil, fmt.Errorf("failed to read private key file: %v", err)
	}

	key, err := keystore.DecryptKey(keyJSON, "")
	if err != nil {
		logger.Error("Failed to decrypt private key", "error", err)
		return nil, fmt.Errorf("failed to decrypt private key: %v", err)
	}

	server := rpc.NewServer()
	submitter := &Submitter{
		server:           server,
		cfg:              cfg,
		ctx:              ctx,
		logger:           logger,
		privateKey:       key,
		serverAddr:       cfg.ServerAddr,
		chainConnections: chainConns,
		taskMap:          sync.Map{},
		nonceMap:         sync.Map{},
	}
	if err := server.Register(submitter); err != nil {
		logger.Error("Failed to register RPC server", "error", err)
		return nil, fmt.Errorf("failed to register RPC server: %v", err)
	}

	submitter.BaseService = *service.NewBaseService(nil, "Submitter", submitter)
	return submitter, nil
}

func (s *Submitter) OnStart() error {
	var err error
	s.listener, err = net.Listen("tcp", s.serverAddr)
	if err != nil {
		s.logger.Error("Failed to start listener", "address", s.serverAddr, "error", err)
		panic(err)
	}

	go s.server.Accept(s.listener)
	return nil
}

func (s *Submitter) OnStop() {
	for chainID, conn := range s.chainConnections {
		conn.ethClient.Close()
		s.logger.Info("Closed connection", "chainID", chainID)
	}
	if s.listener != nil {
		s.listener.Close()
	}
}
