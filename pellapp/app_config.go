package pellapp

import (
	"fmt"

	dvsconfig "github.com/0xPellNetwork/pelldvs/config"

	"intellix/common"
)

type AppConfig struct {
	DvsConfig *dvsconfig.Config `json:"-" mapstructure:"-"`

	RootDir       string `json:"root_dir" mapstructure:"root_dir"`
	GatewayAddr   string `json:"gateway_addr" mapstructure:"gateway_addr"`
	OperatorAddr  string `json:"operator_address" mapstructure:"operator_address"`
	CosmosNodeUri string `json:"cosmos_node_uri" mapstructure:"cosmos_node_uri"`
	CosmosChainId string `json:"cosmos_chain_id" mapstructure:"cosmos_chain_id"`

	WaitBlockCount int64   `json:"wait_block_count" mapstructure:"wait_block_count"`
	GasPrices      string  `json:"gas_prices" mapstructure:"gas_prices"`
	GasAdjustment  float64 `json:"gas_adjustment" mapstructure:"gas_adjustment"`

	PriceTickConverterConfig map[string]map[string]string `json:"price_tick_converter_config" mapstructure:"price_tick_converter_config"`
	ECCKeyPair               common.ECCKeyPair            `json:"ecc_key_pair" mapstructure:"ecc_key_pair"`
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
