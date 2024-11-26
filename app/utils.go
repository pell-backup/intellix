package app

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtxconfig "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"os"
)

func NewClientContext(
	appCodec codec.Codec,
	interfaceRegistry codectypes.InterfaceRegistry,
	txConfigOpts tx.ConfigOptions,
	fromName string, cosmosNodeUrl string, chainId string,
	kr keyring.Keyring,
) client.Context {
	clientCtx := client.Context{}.
		WithCodec(appCodec).
		WithInterfaceRegistry(interfaceRegistry).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithHomeDir(DefaultNodeHome).
		WithViper(Name) // env variable prefix

	// Read the config again to overwrite the default values with the values from the config file
	clientCtx, _ = config.ReadFromClientConfig(clientCtx)

	// textual is enabled by default, we need to re-create the tx config grpc instead of bank keeper.
	txConfigOpts.TextualCoinMetadataQueryFn = authtxconfig.NewGRPCCoinMetadataQueryFn(clientCtx)
	txConfig, err := tx.NewTxConfigWithOptions(clientCtx.Codec, txConfigOpts)
	if err != nil {
		panic(err)
	}
	clientCtx = clientCtx.WithTxConfig(txConfig)

	clientCtx = clientCtx.WithFromName(fromName).WithKeyring(kr)
	c, err := client.NewClientFromNode(cosmosNodeUrl)
	if err != nil {
		panic(err)
	}
	return clientCtx.WithClient(c).WithNodeURI(cosmosNodeUrl).WithChainID(chainId)
}
