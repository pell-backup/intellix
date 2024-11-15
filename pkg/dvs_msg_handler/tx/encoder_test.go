package tx

import (
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
		Task:      &types.TaskRequest{},
		PriceFeed: &types.PriceFeedParam{},
	})
	require.NoError(t, err)

	txBz, err := coder.Encode(builder.GetTx())
	require.NoError(t, err)

	// decode before register: should has error
	//_, err = coder.Decode(txBz)
	//require.EqualError(t, err, "unable to resolve type URL /intellix.price.ProcessRequestPriceFeedIn: tx parse error")

	// decode after register
	types.RegisterInterfaces(registry)
	_, err = coder.Decode(txBz)
	require.NoError(t, err)
}

func TestTxMsgs(t *testing.T) {
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	coder := NewDefaultDecoder(cdc)

	builder := NewBuilder(cdc)
	err := builder.SetMsgs(&types.ProcessRequestPriceFeedIn{
		Task:      nil,
		PriceFeed: nil,
	})
	require.NoError(t, err)

	for _, msg := range builder.GetTx().GetMsgs() {
		url := sdk.MsgTypeURL(msg)
		require.NotEmpty(t, url)
	}

	txBz, err := coder.Encode(builder.GetTx())
	require.NoError(t, err)

	types.RegisterInterfaces(registry)
	decodedTx, err := coder.Decode(txBz)
	require.NoError(t, err)

	require.Greater(t, len(decodedTx.GetMsgs()), 0)

	for _, msg := range decodedTx.GetMsgs() {
		url := sdk.MsgTypeURL(msg)
		require.NotEmpty(t, url)
	}
}
