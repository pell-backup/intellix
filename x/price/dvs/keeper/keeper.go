package keeper

import (
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"fmt"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/spf13/pflag"
	"intellix/x/price/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger
		clientCtx    client.Context

		operatorAddress string
		gasPrices       string
		gasAdjustment   float64

		waitBlockCount int64 // price feed wait block count
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	clientCtx client.Context,

	operatorAddress string,
	waitBlockCount int64,

	gasPrices string,
	gasAdjustment float64,
) Keeper {
	if gasPrices == "" {
		gasPrices = "0.1uatom"
	}
	if gasAdjustment == 0 {
		gasAdjustment = 1.5
	}
	if waitBlockCount == 0 {
		panic("waitBlockCount can't be nil or zero")
	}

	k := Keeper{
		cdc:          cdc,
		storeService: storeService,
		logger:       logger,
		clientCtx:    clientCtx,

		operatorAddress: operatorAddress,
		waitBlockCount:  waitBlockCount,
		gasPrices:       gasPrices,
		gasAdjustment:   gasAdjustment,
	}

	k.SetOperatorAddress(operatorAddress)
	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetOperatorAddress(ctx sdk.Context) string {
	if k.operatorAddress == "" {
		panic("Operator address not set")
	}
	return k.operatorAddress
}

func (k *Keeper) SetOperatorAddress(address string) {
	if err := sdk.VerifyAddressFormat(sdk.AccAddress(address)); err != nil {
		panic(err)
	}
	k.operatorAddress = address
}

func (k Keeper) GetLatestBlock(ctx sdk.Context) (*cmttypes.Block, error) {
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
func (k Keeper) SignAndBroadcastTx(ctx sdk.Context, msg sdk.Msg) error {
	txf, err := k.prepareTxFactory(ctx)
	if err != nil {
		return fmt.Errorf("failed to prepare tx factory: %w", err)
	}

	txBuilder, err := txf.BuildUnsignedTx(msg)
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

// prepareTxFactory prepare tx factory
func (k Keeper) prepareTxFactory(ctx sdk.Context) (tx.Factory, error) {
	txf, err := tx.NewFactoryCLI(k.clientCtx, &pflag.FlagSet{})
	if err != nil {
		return tx.Factory{}, err
	}

	txf = txf.WithGasPrices(k.gasPrices).WithGasAdjustment(k.gasAdjustment)
	txf = txf.WithChainID(ctx.ChainID())
	txf = txf.WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	return txf, nil
}
