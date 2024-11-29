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
