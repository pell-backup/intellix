package taskgateway

import (
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/0xPellNetwork/pelldvs/crypto/bls"
	contractdataoracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/ethereum/go-ethereum/common"
)

type ChainConfig struct {
	EthEndpoint     string `mapstructure:"eth_endpoint"`
	ContractAddress string `mapstructure:"contract_address"`
	ChainID         int64  `mapstructure:"chain_id"`
	GasLimit        uint64 `mapstructure:"gas_limit"`
}

type TaskGatewayCfg struct {
	ServerAddr          string                `mapstructure:"server_addr"`
	SenderAddress       string                `mapstructure:"sender_address"`
	PrivateKeyStorePath string                `mapstructure:"private_key_store_path"`
	Chains              map[int64]ChainConfig `mapstructure:"chains"`
}

func (t TaskGatewayCfg) Validate() error {
	if t.SenderAddress == "" {
		return fmt.Errorf("sender address cannot be empty")
	}
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
		if cfg.EthEndpoint == "" {
			return fmt.Errorf("eth endpoint for chain %d cannot be empty", chainID)
		}
		if cfg.ContractAddress == "" {
			return fmt.Errorf("contract address for chain %d cannot be empty", chainID)
		}
	}
	return nil
}

func convertAddressToString(addrStr string) (*common.Address, error) {
	mixedCaseAddress, err := common.NewMixedcaseAddressFromString(addrStr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert address: %w", err)
	}
	addr := mixedCaseAddress.Address()
	return &addr, nil
}

func convertNonSignersPubkeysG1(pb [][]byte) []contractdataoracle.BN254G1Point {
	list := make([]contractdataoracle.BN254G1Point, len(pb))
	for i, p := range pb {
		if len(p) < 64 {
			continue // Skip invalid points
		}
		list[i] = contractdataoracle.BN254G1Point{
			X: new(big.Int).SetBytes(p[:32]),
			Y: new(big.Int).SetBytes(p[32:]),
		}
	}
	return list
}

func convertToBN254G1Point(input *bls.G1Point) contractdataoracle.BN254G1Point {
	if input == nil {
		return contractdataoracle.BN254G1Point{
			X: new(big.Int),
			Y: new(big.Int),
		}
	}
	output := contractdataoracle.BN254G1Point{
		X: input.X.BigInt(new(big.Int)),
		Y: input.Y.BigInt(new(big.Int)),
	}
	return output
}

func convertQuorumApks(pb [][]byte) []contractdataoracle.BN254G1Point {
	list := make([]contractdataoracle.BN254G1Point, 0, len(pb))
	for _, apk := range pb {
		if len(apk) == 0 {
			continue // Skip empty APKs
		}
		tapk := bls.NewZeroG1Point()
		if err := tapk.Unmarshal(apk); err != nil {
			continue // Skip invalid points
		}
		list = append(list, convertToBN254G1Point(tapk))
	}
	return list
}

func convertApkG2(pb []byte) contractdataoracle.BN254G2Point {
	if len(pb) < 128 {
		return contractdataoracle.BN254G2Point{}
	}

	return contractdataoracle.BN254G2Point{
		X: [2]*big.Int{
			new(big.Int).SetBytes(pb[:32]),
			new(big.Int).SetBytes(pb[32:64]),
		},
		Y: [2]*big.Int{
			new(big.Int).SetBytes(pb[64:96]),
			new(big.Int).SetBytes(pb[96:]),
		},
	}
}

func convertSigma(pb []byte) contractdataoracle.BN254G1Point {
	if len(pb) < 64 {
		return contractdataoracle.BN254G1Point{}
	}

	return contractdataoracle.BN254G1Point{
		X: new(big.Int).SetBytes(pb[:32]),
		Y: new(big.Int).SetBytes(pb[32:]),
	}
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

func validateBLSComponents(data *RPCValidatedData) error {
	if data == nil {
		return fmt.Errorf("validated data is nil")
	}

	// Validate SignersApkG2 and SignersAggSigG1
	if len(data.SignersApkG2) < 128 {
		return fmt.Errorf("invalid SignersApkG2 length: got %d, want >= 128", len(data.SignersApkG2))
	}
	if len(data.SignersAggSigG1) < 64 {
		return fmt.Errorf("invalid SignersAggSigG1 length: got %d, want >= 64", len(data.SignersAggSigG1))
	}

	// Validate QuorumApks
	if len(data.QuorumApksG1) == 0 {
		return fmt.Errorf("no QuorumApks provided")
	}
	if len(data.QuorumApksG1) != len(data.QuorumApkIndices) {
		return fmt.Errorf("QuorumApks length mismatch: got %d APKs but %d indices",
			len(data.QuorumApksG1), len(data.QuorumApkIndices))
	}

	// Validate indices
	if len(data.TotalStakeIndices) == 0 {
		return fmt.Errorf("no TotalStakeIndices provided")
	}

	return nil
}

const (
	TaskTypePriceFeed int64 = 1
	TaskTypeProcessor int64 = 3
)
