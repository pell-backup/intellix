package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
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

type ECCKeyPair struct {
	ECCPrivateKey string `json:"eccprivate_key"`
	ECCPublicKey  string `json:"ecc_public_key"`
}

type Config struct {
	ServerAddr          string                 `json:"server_addr"`
	PrivateKeyStorePath string                 `json:"private_key_store_path"`
	Chains              map[uint64]ChainConfig `json:"chains"`
	ECCKeyPair          ECCKeyPair             `json:"ecc_key_pair"`
}

func LoadConfig(cfgPath string) (*Config, error) {
	var cfg Config
	input, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(input, &cfg)
	return &cfg, err
}

func (t Config) Validate() error {
	if t.ServerAddr == "" {
		return fmt.Errorf("server address cannot be empty")
	}
	if t.PrivateKeyStorePath == "" {
		return fmt.Errorf("key_store_path cannot be empty")
	}
	if _, err := os.Stat(t.PrivateKeyStorePath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("key_store_path does not exist")
	}
	if len(t.Chains) == 0 {
		return fmt.Errorf("no chain configs provided")
	}
	for chainID, cfg := range t.Chains {
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("chain %d: %w", chainID, err)
		}
	}
	return nil
}
