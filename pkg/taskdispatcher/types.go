package taskdispatcher

import (
	"fmt"
)

type ChainConfig struct {
	ChainID         string `mapstructure:"chain_id"`
	EthURL          string `mapstructure:"eth_url"`
	ContractAddress string `mapstructure:"contract_address"`
}

// Validate checks if the ChainConfig is valid
func (c ChainConfig) Validate() error {
	if c.ChainID == "" {
		return fmt.Errorf("chain_id cannot be empty")
	}
	if c.EthURL == "" {
		return fmt.Errorf("eth_url cannot be empty")
	}
	if c.ContractAddress == "" {
		return fmt.Errorf("contract_address cannot be empty")
	}
	return nil
}
