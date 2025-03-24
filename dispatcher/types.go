package taskdispatcher

import (
	"bytes"
	"fmt"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/spf13/viper"

	"intellix/config"
)

type Config struct {
	Chains map[uint64]config.ChainConfig `json:"chains"`
}

func LoadConfig(filepath string) (*Config, error) {
	var cfg Config
	vp := viper.New()
	vp.SetConfigFile(filepath)
	if err := vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	if err := vp.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &cfg, nil
}

func (c Config) Validate() error {
	if len(c.Chains) == 0 {
		return fmt.Errorf("no chain specified")
	}
	for _, chain := range c.Chains {
		if err := chain.Validate(); err != nil {
			return err
		}
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
	// decode cbor
	reader := bytes.NewReader(data)
	decoder := cbor.NewDecoder(reader)

	var key string
	result := make(map[string]interface{})
	for i := 0; ; i++ {
		var item interface{}
		err := decoder.Decode(&item)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}

		if i%2 == 0 {
			// key
			if str, ok := item.(string); ok {
				key = str
			}
		} else {
			// value
			result[key] = item
		}
	}

	var sp = &ScriptParam{}
	if scriptId, ok := result["scriptId"]; ok {
		sp.ScriptId = scriptId.(uint64)
	}
	if params, ok := result["params"]; ok {
		sp.Params = params.([]byte)
	}

	return sp, nil
}

const (
	TaskTypePrice  int64 = 1
	TaskTypeScript int64 = 3
)
