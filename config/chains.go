package config

import (
	"fmt"
)

type ChainConfig struct {
	RPCURL          string `json:"rpc_url"`
	WSURL           string `json:"ws_url"`
	ContractAddress string `json:"contract_address"`
	ChainID         uint64 `json:"chain_id"`
	GasLimit        uint64 `json:"gas_limit"`
}

func (c ChainConfig) Validate() error {
	if c.ChainID == 0 {
		return fmt.Errorf("chain_id cannot be empty")
	}
	if c.WSURL == "" {
		return fmt.Errorf("ws_url cannot be empty")
	}
	if c.RPCURL == "" {
		return fmt.Errorf("rpc_url cannot be empty")
	}
	if c.ContractAddress == "" {
		return fmt.Errorf("contract_address cannot be empty")
	}
	return nil
}
