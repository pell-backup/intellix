package server

import (
	"context"
	"fmt"
	taskgateway "intellix/gateway"
	sdktypes "intellix/sdk/types"

	"github.com/0xPellNetwork/pelldvs/crypto/bls"
	"github.com/0xPellNetwork/pelldvs/libs/log"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/spf13/pflag"
)

type Server struct {
	logger        log.Logger
	clientCtx     client.Context
	cosmosChainId string
	key           *keyring.Record

	taskGatewayClient *taskgateway.Client

	wsEndpoint      string
	operatorAddress string
	gasPrices       string
	gasAdjustment   float64
	waitBlockCount  int64 // price feed wait block count
	blsKeyPair      *bls.KeyPair
}

func NewServer(
	logger log.Logger,
	clientCtx client.Context,
	key *keyring.Record,
	cosmosChainId string,

	wsEndpoint string,
	gatewayAddr string,
	operatorAddress string,
	waitBlockCount int64,
	blsKeyPath, blsKeyPassword string,

	gasPrices string,
	gasAdjustment float64,
) (Server, error) {
	if gasPrices == "" {
		gasPrices = "1stake"
	}
	if gasAdjustment == 0 {
		gasAdjustment = 1.5
	}
	if waitBlockCount == 0 {
		waitBlockCount = 1
	}

	k := Server{
		logger:        logger,
		clientCtx:     clientCtx,
		key:           key,
		cosmosChainId: cosmosChainId,

		wsEndpoint:      wsEndpoint,
		operatorAddress: operatorAddress,
		waitBlockCount:  waitBlockCount,
		gasPrices:       gasPrices,
		gasAdjustment:   gasAdjustment,
	}

	if blsKeyPath != "" && blsKeyPassword != "" {
		var err error
		k.blsKeyPair, err = bls.ReadPrivateKeyFromFile(blsKeyPath, blsKeyPassword)
		if err != nil {
			return Server{}, fmt.Errorf("failed to load BLS key pair: %w", err)
		}
	}

	if operatorAddress != "" {
		k.SetOperatorAddress(operatorAddress)
		_, err := clientCtx.Keyring.Key(clientCtx.GetFromName())
		if err != nil {
			return Server{}, fmt.Errorf("operator key not found in keyring: %w", err)
		}
	}

	taskGatewayClient, err := taskgateway.NewClient(gatewayAddr, logger)
	if err != nil {
		return Server{}, err
	}
	k.taskGatewayClient = taskGatewayClient

	return k, nil
}

func (k *Server) SetOperatorAddress(address string) {
	if err := sdk.VerifyAddressFormat(sdk.AccAddress(address)); err != nil {
		panic(err)
	}
	k.operatorAddress = address
}

// SignAndBroadcastTx signs and broadcasts a transaction
func (k *Server) SignAndBroadcastTx(ctx sdktypes.Context, msg sdk.Msg) error {
	txf, err := k.prepareTxFactory(ctx)
	if err != nil {
		return fmt.Errorf("failed to prepare tx factory: %w", err)
	}

	address, err := k.key.GetAddress()
	if err != nil {
		return err
	}
	execMsg := authz.NewMsgExec(address, []sdk.Msg{msg})

	txBuilder, err := txf.BuildUnsignedTx(&execMsg)
	if err != nil {
		return fmt.Errorf("failed to build unsigned tx: %w", err)
	}

	err = tx.Sign(ctx, txf, k.clientCtx.GetFromName(), txBuilder, true)
	if err != nil {
		return fmt.Errorf("failed to sign tx: %w", err)
	}

	txBytes, err := k.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return fmt.Errorf("failed to encode tx: %w", err)
	}

	res, err := k.clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return fmt.Errorf("failed to broadcast tx: %w", err)
	}

	if res.Code != 0 {
		return fmt.Errorf("tx failed with code %d: %s", res.Code, res.RawLog)
	}

	return nil
}

func (k *Server) SenderAddress() (sdk.AccAddress, error) {
	addr, err := k.key.GetAddress()
	if err != nil {
		return sdk.AccAddress{}, fmt.Errorf("failed to get address: %w", err)
	}
	return addr, nil
}

// prepareTxFactory prepare tx factory
func (k *Server) prepareTxFactory(ctx sdktypes.Context) (tx.Factory, error) {
	txf, err := tx.NewFactoryCLI(k.clientCtx, &pflag.FlagSet{})
	if err != nil {
		return tx.Factory{}, err
	}

	txf = txf.WithGasPrices(k.gasPrices).WithGasAdjustment(k.gasAdjustment)
	txf = txf.WithChainID(k.cosmosChainId)
	txf = txf.WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	addr, err := k.SenderAddress()
	if err != nil {
		return tx.Factory{}, err
	}

	// 获取账户信息
	accRetriever := authtypes.AccountRetriever{}
	acc, err := accRetriever.GetAccount(k.clientCtx, addr)
	if err != nil {
		return tx.Factory{}, err
	}
	txf = txf.WithAccountNumber(acc.GetAccountNumber()).WithSequence(acc.GetSequence())

	return txf, nil
}

func (k *Server) GetLatestBlock(ctx context.Context) (*cmttypes.Block, error) {
	node, err := k.clientCtx.GetNode()
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	// get last block height
	latestHeight, err := node.Block(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block height: %w", err)
	}

	// get last block info by height
	block, err := node.Block(ctx, &latestHeight.Block.Height)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}

	return block.Block, nil
}
