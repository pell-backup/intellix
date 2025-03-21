package dispatcher

import (
	"bytes"
	"fmt"
	"io"

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

type VRFTaskParam struct {
	NumWords uint64
}

func ParseVRFTaskParam(data []byte) (*VRFTaskParam, error) {
	decoder := cbor.NewDecoder(bytes.NewReader(data))
	result := make(map[string]interface{})

	// Continuously decode key-value pairs until EOF
	for {
		var key interface{}
		var val interface{}

		// Decode the key
		if err := decoder.Decode(&key); err != nil {
			if err == io.EOF {
				// We've reached the end, break out of the loop
				break
			}
			return nil, fmt.Errorf("failed to decode key: %w", err)
		}

		// Decode the value
		if err := decoder.Decode(&val); err != nil {
			if err == io.EOF {
				// End of data
				break
			}
			return nil, fmt.Errorf("failed to decode value: %w", err)
		}

		// Convert key to string
		strKey, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("expected key to be string, got %T", key)
		}

		// Store in map
		result[strKey] = val
	}

	// Once the loop is done, extract the field you need
	var out VRFTaskParam
	if val, ok := result["numWords"]; ok {
		numWords, ok := val.(uint64)
		if !ok {
			return nil, fmt.Errorf("expected numWords to be uint64, got %T", val)
		}
		out.NumWords = numWords
	} else {
		// If "numWords" is required but missing, return an error
		return nil, fmt.Errorf("missing numWords in data")
	}

	return &out, nil
}
