package gateway

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

type Config struct {
	ServerAddr          string                 `json:"server_addr"`
	PrivateKeyStorePath string                 `json:"private_key_store_path"`
	Chains              map[uint64]ChainConfig `json:"chains"`
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

type RespondToTaskResponse struct {
	Error string `json:"error"`
}

// RPCTaskRaw represents a serializable version of TaskRaw
type RPCTaskRaw struct {
	TaskType                  int64  `json:"task_type"`
	TaskIndex                 uint32 `json:"task_index"`
	RequestID                 []byte `json:"request_id"`
	FeeToken                  string `json:"fee_token"`
	Payment                   string `json:"payment"`
	RequestData               []byte `json:"request_data"`
	CallbackAddress           string `json:"callback_address"`
	CallbackFunctionID        []byte `json:"callback_function_id"`
	TaskCreatedBlock          uint32 `json:"task_created_block"`
	QuorumNumbers             []byte `json:"quorum_numbers"`
	QuorumThresholdPercentage uint32 `json:"quorum_threshold_percentage"`
	AdvanceDecode             bool   `json:"advance_decode"`
}

type RPCTaskResponse struct {
	ReferenceTaskIndex uint32 `json:"reference_task_index"`
	Data               []byte `json:"data"`
}

type RPCValidatedData struct {
	Data                         []byte     `json:"data,omitempty"`
	Error                        string     `json:"error,omitempty"`
	Hash                         []byte     `json:"hash,omitempty"`
	NonSignersPubkeysG1          [][]byte   `json:"non_signers_pubkeys_g_1,omitempty"`
	QuorumApksG1                 [][]byte   `json:"quorum_apks_g_1,omitempty"`
	SignersApkG2                 []byte     `json:"signers_apk_g_2,omitempty"`
	SignersAggSigG1              []byte     `json:"signers_agg_sig_g_1,omitempty"`
	NonSignerQuorumBitmapIndices []uint32   `json:"non_signer_quorum_bitmap_indices,omitempty"`
	QuorumApkIndices             []uint32   `json:"quorum_apk_indices,omitempty"`
	TotalStakeIndices            []uint32   `json:"total_stake_indices,omitempty"`
	NonSignerStakeIndices        [][]uint32 `json:"non_signer_stake_indices,omitempty"`
}

type RPCVoteFinalizedRequestIn struct {
	ChainID        int64             `json:"chain_id"`
	TaskRaw        *RPCTaskRaw       `json:"task_raw"`
	ValidatedData  *RPCValidatedData `json:"validated_data"`
	RespToTaskData []byte            `json:"resp_to_task_data"` // decoded data
}

const (
	TaskTypePriceFeed int64 = 1
	TaskTypeProcessor int64 = 3
)
