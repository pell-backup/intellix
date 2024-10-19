package dvsservermanager

import (
	"cosmossdk.io/collections/colltest"
	cosmossdk_io_math "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
	"intellix/pkg/dvs_msg_handler/tx"
	"intellix/proto/intellix/x/price/types"
	"intellix/x/price/dvs/server"
	dvstypes "intellix/x/price/dvs/types"
	"testing"
)

func genMsgData(t *testing.T, cdc codec.Codec) []byte {
	coder := tx.NewDefaultDecoder(cdc)

	builder := tx.NewBuilder(cdc)
	err := builder.SetMsgs(&types.ProcessPriceFeedMsg{
		Height:  1,
		ChainId: cosmossdk_io_math.Int{},
	})
	require.NoError(t, err)

	txBz, err := coder.Encode(builder.GetTx())
	require.NoError(t, err)

	_, err = coder.Decode(txBz)
	require.NoError(t, err)

	return txBz
}

func TestProcessRequestHandler(t *testing.T) {
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	handler := NewProcessRequestHandler(tx.NewDefaultDecoder(cdc))
	mockStore, _ := colltest.MockStore()
	k := server.NewKeeper(cdc, mockStore, nil, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	// register interface
	dvstypes.RegisterInterfaces(registry)

	// register msg router
	dvstypes.RegisterDvsProcessRequestServer(handler, server.NewDvsProcessRequestServer(k))

	// get handler
	_, ok := handler.(*ProcessRequestHandler)
	require.Equal(t, ok, true)

}
