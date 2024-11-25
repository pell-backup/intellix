package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/x/authz"
	pkgcontext "intellix/pkg/context"
	"intellix/pkg/taskgateway"
	"intellix/x/price/types"

	"github.com/0xPellNetwork/pelldvs/libs/log"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/spf13/pflag"
)

type (
	Server struct {
		logger    log.Logger
		clientCtx client.Context
		key       *keyring.Record

		operatorAddress string
		gasPrices       string
		gasAdjustment   float64
		waitBlockCount  int64 // price feed wait block count

		taskGatewayClient *taskgateway.Client
	}
)

func NewServer(
	logger log.Logger,
	clientCtx client.Context,
	key *keyring.Record,

	gatewayAddr string,
	operatorAddress string,
	waitBlockCount int64,

	gasPrices string,
	gasAdjustment float64,
) (Server, error) {
	if gasPrices == "" {
		gasPrices = "0.1uatom"
	}
	if gasAdjustment == 0 {
		gasAdjustment = 1.5
	}
	if waitBlockCount == 0 {
		waitBlockCount = 10
	}

	k := Server{
		logger:    logger,
		clientCtx: clientCtx,
		key:       key,

		operatorAddress: operatorAddress,
		waitBlockCount:  waitBlockCount,
		gasPrices:       gasPrices,
		gasAdjustment:   gasAdjustment,
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

// Logger returns a module-specific logger.
func (k *Server) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Server) GetOperatorAddress(ctx context.Context) string {
	if k.operatorAddress == "" {
		panic("Operator address not set")
	}
	return k.operatorAddress
}

func (k *Server) SetOperatorAddress(address string) {
	if err := sdk.VerifyAddressFormat(sdk.AccAddress(address)); err != nil {
		panic(err)
	}
	k.operatorAddress = address
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

// SignAndBroadcastTx signs and broadcasts a transaction
func (k *Server) SignAndBroadcastTx(ctx pkgcontext.Context, msgs []sdk.Msg) error {
	txf, err := k.prepareTxFactory(ctx)
	if err != nil {
		return fmt.Errorf("failed to prepare tx factory: %w", err)
	}

	address, err := k.key.GetAddress()
	if err != nil {
		return err
	}
	execMsg := authz.NewMsgExec(address, msgs)

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

	// base64 encode tx bytes
	base64Tx := base64.StdEncoding.EncodeToString(txBytes)
	k.Logger().Info("Broadcasting tx", "tx", base64Tx)

	res, err := k.clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return fmt.Errorf("failed to broadcast tx: %w", err)
	}

	if res.Code != 0 {
		return fmt.Errorf("tx failed with code %d: %s", res.Code, res.RawLog)
	}

	return nil
}

// prepareTxFactory prepare tx factory
func (k *Server) prepareTxFactory(ctx pkgcontext.Context) (tx.Factory, error) {
	txf, err := tx.NewFactoryCLI(k.clientCtx, &pflag.FlagSet{})
	if err != nil {
		return tx.Factory{}, err
	}

	txf = txf.WithGasPrices(k.gasPrices).WithGasAdjustment(k.gasAdjustment)
	txf = txf.WithChainID(fmt.Sprintf("%d", ctx.ChainID()))
	txf = txf.WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	return txf, nil
}
