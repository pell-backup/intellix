package taskdispatcher

import (
	"bytes"
	"fmt"
	"math/big"

	cbor "github.com/fxamacker/cbor/v2"
)

type Config struct {
	Chains     []*ChainConfig `mapstructure:"chains"`
	DvsAddress string         `mapstructure:"dvs_address"`
}

func (c Config) Validate() error {
	if len(c.Chains) == 0 {
		return fmt.Errorf("no chain specified")
	}
	if c.DvsAddress == "" {
		return fmt.Errorf("dvs address is required")
	}
	for _, chain := range c.Chains {
		if err := chain.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type ChainConfig struct {
	ChainID         uint64 `mapstructure:"chain_id"`
	EthURL          string `mapstructure:"eth_url"`
	ContractAddress string `mapstructure:"contract_address"`
}

// Validate checks if the ChainConfig is valid
func (c ChainConfig) Validate() error {
	if c.ChainID == 0 {
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

type PriceFeedParam struct {
	BaseSymbol  string
	QuoteSymbol string
}

func ParsePriceFeed(data []byte) (*PriceFeedParam, error) {
	decoder := cbor.NewDecoder(bytes.NewReader(data[:]))
	var pDecoded1 interface{}
	taskDataTmp := make([]interface{}, 0)
	for i := 0; i < 6; i++ {
		err := decoder.Decode(&pDecoded1)
		if err != nil {
			return nil, err
		}

		taskDataTmp = append(taskDataTmp, pDecoded1)
	}

	baseSymbol := taskDataTmp[1].(string)
	quoteSymbol := taskDataTmp[3].(string)

	return &PriceFeedParam{baseSymbol, quoteSymbol}, nil
}

type ScriptParam struct {
	ScriptId uint64
	Params   []byte
}

func ParseScript(data []byte) (*ScriptParam, error) {
	cborData := data[24:]

	var rawMap map[string]interface{}
	if err := cbor.Unmarshal(cborData, &rawMap); err != nil {
		return nil, fmt.Errorf("CBOR decode failed: %v", err)
	}

	req := &ScriptParam{}

	// scriptId
	if scriptIdRaw, ok := rawMap["scriptId"]; ok {
		switch v := scriptIdRaw.(type) {
		case uint64:
			req.ScriptId = v
		case *big.Int:
			if v.IsUint64() {
				req.ScriptId = v.Uint64()
			} else {
				return nil, fmt.Errorf("scriptId exceeds uint64 range")
			}
		default:
			return nil, fmt.Errorf("invalid scriptId type")
		}
	}

	// params
	if paramsRaw, ok := rawMap["params"]; ok {
		if params, ok := paramsRaw.([]byte); ok {
			req.Params = params
		} else {
			return nil, fmt.Errorf("invalid params type")
		}
	}

	return req, nil
}

const (
	TaskTypePrice int64 = 1
)
