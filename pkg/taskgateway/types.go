package taskgateway

import (
	"errors"
	"fmt"
	"github.com/0xPellNetwork/pelldvs/crypto/bls"
	dataOracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"os"
)

type TaskGatewayCfg struct {
	ServerAddr          string `mapstructure:"server_addr"`
	EthEndpoint         string `mapstructure:"eth_endpoint"`
	SenderAddress       string `mapstructure:"sender_address"`
	ContractAddress     string `mapstructure:"contract_address"`
	PrivateKeyStorePath string `mapstructure:"private_key_store_path"`
}

func (t TaskGatewayCfg) Validate() error {
	if t.EthEndpoint == "" {
		return fmt.Errorf("eth endpoint cannot be empty")
	}
	if t.ContractAddress == "" {
		return fmt.Errorf("contract address cannot be empty")
	}
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

func convertNonSignersPubkeysG1(pb [][]byte) []dataOracle.BN254G1Point {
	list := []dataOracle.BN254G1Point{}
	for i, p := range pb {
		list[i] = dataOracle.BN254G1Point{
			X: new(big.Int).SetBytes(p[:32]),
			Y: new(big.Int).SetBytes(p[32:]),
		}
	}
	return list
}

func convertToBN254G1Point(input *bls.G1Point) dataOracle.BN254G1Point {
	output := dataOracle.BN254G1Point{
		X: input.X.BigInt(big.NewInt(0)),
		Y: input.Y.BigInt(big.NewInt(0)),
	}
	return output
}

func convertQuorumApks(pb [][]byte) []dataOracle.BN254G1Point {
	list := []dataOracle.BN254G1Point{}
	for _, apk := range pb {
		tapk := bls.NewZeroG1Point()
		_ = tapk.Unmarshal(apk)
		list = append(list, convertToBN254G1Point(tapk))
	}
	return list
}

func convertApkG2(pb []byte) *dataOracle.BN254G2Point {
	return &dataOracle.BN254G2Point{
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

func convertSigma(pb []byte) *dataOracle.BN254G1Point {
	return &dataOracle.BN254G1Point{
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
}

// RPCPriceFeedResponse represents a serializable version of PriceFeedResponse
type RPCPriceFeedResponse struct {
	ReferenceTaskIndex uint32 `json:"reference_task_index"`
	Price              string `json:"price"`
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

// RPCVoteFinalizedRequestPrice represents a serializable version of MsgVoteFinalizedRequestPrice
type RPCVoteFinalizedRequestPrice struct {
	ChainID           int64                 `json:"chain_id"`
	TaskRaw           *RPCTaskRaw           `json:"task_raw"`
	ValidatedData     *RPCValidatedData     `json:"validated_data"`
	PriceFeedResponse *RPCPriceFeedResponse `json:"price_feed_response"`
}
