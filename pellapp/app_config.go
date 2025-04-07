package pellapp

import (
	"fmt"

	dvsconfig "github.com/0xPellNetwork/pelldvs/config"
)

type AppConfig struct {
	DvsConfig *dvsconfig.Config `mapstructure:"-"`

	RootDir       string `mapstructure:"root_dir"`
	GatewayAddr   string `mapstructure:"gateway_addr"`
	OperatorAddr  string `mapstructure:"operator_address"`
	CosmosNodeUri string `mapstructure:"cosmos_node_uri"`
	CosmosChainId string `mapstructure:"cosmos_chain_id"`

	WaitBlockCount int64   `mapstructure:"wait_block_count"`
	GasPrices      string  `mapstructure:"gas_prices"`
	GasAdjustment  float64 `mapstructure:"gas_adjustment"`

	PriceTickConverterConfig map[string]map[string]string `mapstructure:"price_tick_converter_config"`
	ApiKey                   map[string]string            `mapstructure:"api_key"`
	SymbolList               map[string]map[string]string `mapstructure:"symbol_list"`
}

func (p AppConfig) Validate() error {
	if p.OperatorAddr == "" {
		return fmt.Errorf("no operator address provided")
	}
	if p.DvsConfig == nil || p.DvsConfig.ValidateBasic() != nil {
		return fmt.Errorf("invalid pell config")
	}
	if p.GatewayAddr == "" {
		return fmt.Errorf("no gateway address provided")
	}
	if p.CosmosNodeUri == "" {
		return fmt.Errorf("no cosmos node uri provided")
	}
	return nil
}
