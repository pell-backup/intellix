package tx

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"intellix/x/price/dvs/types"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	coder := NewDefaultDecoder(cdc)

	builder := NewBuilder(cdc)
	err := builder.SetMsgs(&types.ProcessRequestPriceFeedIn{
		TaskIndex: 1,
	})
	require.NoError(t, err)

	txBz, err := coder.Encode(builder.GetTx())
	require.NoError(t, err)

	// decode before register: should has error
	_, err = coder.Decode(txBz)
	require.EqualError(t, err, "unable to resolve type URL /intellix.price.ProcessRequestPriceFeedIn: tx parse error")

	// decode after register
	types.RegisterInterfaces(registry)
	decodedTx, err := coder.Decode(txBz)
	require.NoError(t, err)
	_tx, ok := decodedTx.(sdk.Tx)
	if !ok {
		require.Error(t, fmt.Errorf("unable to cast to sdk.Tx"))
	}
	_ = _tx
}

func TestTxMsgs(t *testing.T) {
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	coder := NewDefaultDecoder(cdc)

	builder := NewBuilder(cdc)
	err := builder.SetMsgs(&types.ProcessRequestPriceFeedIn{
		TaskIndex: 1,
	})
	require.NoError(t, err)

	for _, msg := range builder.GetTx().GetMsgs() {
		url := sdk.MsgTypeURL(msg)
		require.Equal(t, url, "/intellix.price.ProcessRequestPriceFeedIn")
	}

	txBz, err := coder.Encode(builder.GetTx())
	require.NoError(t, err)

	types.RegisterInterfaces(registry)
	decodedTx, err := coder.Decode(txBz)
	require.NoError(t, err)

	_tx, ok := decodedTx.(sdk.Tx)
	if !ok {
		require.Error(t, fmt.Errorf("unable to cast to sdk.Tx"))
	}
	require.Greater(t, len(_tx.GetMsgs()), 0)

	for _, msg := range _tx.GetMsgs() {
		url := sdk.MsgTypeURL(msg)
		require.Equal(t, url, "/intellix.price.ProcessPriceFeedMsg")
	}
}
