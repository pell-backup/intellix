package taskdispatcher

import (
	"bytes"
	cbor "github.com/fxamacker/cbor/v2"
)

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
