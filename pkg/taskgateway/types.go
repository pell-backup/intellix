package taskgateway

import (
	"cosmossdk.io/math"
	"fmt"
	"github.com/0xPellNetwork/pelldvs/crypto/bls"
	dataOracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/common"
	"intellix/pkg/pelldvs/types"
	pricetypes "intellix/x/price/types"
	"math/big"
)

type TaskGatewayCfg struct {
	ServerAddr      string `mapstructure:"server_addr"`
	EthEndpoint     string `mapstructure:"eth_endpoint"`
	SenderAddress   string `mapstructure:"sender_address"`
	ContractAddress string `mapstructure:"contract_address"`
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
	list := make([]dataOracle.BN254G1Point, len(pb))
	for i, p := range pb {
		list[i] = dataOracle.BN254G1Point{
			X: big.NewInt(0).SetBytes(p[:32]),
			Y: big.NewInt(0).SetBytes(p[32:64]),
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
	list := make([]dataOracle.BN254G1Point, len(pb))
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
			big.NewInt(0).SetBytes(pb[:32]),
			big.NewInt(0).SetBytes(pb[32:64]),
		},
		Y: [2]*big.Int{
			big.NewInt(0).SetBytes(pb[64:96]),
			big.NewInt(0).SetBytes(pb[96:]),
		},
	}
}

func convertSigma(pb []byte) *dataOracle.BN254G1Point {
	return &dataOracle.BN254G1Point{
		X: new(big.Int).SetBytes(pb[:32]),
		Y: new(big.Int).SetBytes(pb[32:]),
	}
}

func convertNonSignerStakeIndices(list []*types.NonSignerStakeIndice) [][]uint32 {
	slice := make([][]uint32, len(list))
	for i, l := range list {
		slice[i] = l.NonSignerStakeIndice
	}
	return slice
}

type RespondToTaskResponse struct {
	Error string `json:"error"`
}

// RPCTaskRaw represents a serializable version of TaskRaw
type RPCTaskRaw struct {
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

// RPCVoteFinalizedRequestPrice represents a serializable version of MsgVoteFinalizedRequestPrice
type RPCVoteFinalizedRequestPrice struct {
	TaskRaw           *RPCTaskRaw           `json:"task_raw"`
	ValidatedData     []byte                `json:"validated_data"`
	PriceFeedResponse *RPCPriceFeedResponse `json:"price_feed_response"`
}

// ToProtoMessage converts RPCVoteFinalizedRequestPrice to MsgVoteFinalizedRequestPrice
func (r *RPCVoteFinalizedRequestPrice) ToProtoMessage() (*pricetypes.MsgVoteFinalizedRequestPrice, error) {
	if r == nil {
		return nil, fmt.Errorf("nil RPCVoteFinalizedRequestPrice")
	}
	payment, ok := math.NewIntFromString(r.TaskRaw.Payment)
	if !ok {
		return nil, fmt.Errorf("failed to convert payment: %s", r.TaskRaw.Payment)
	}

	taskRaw := &pricetypes.TaskRaw{
		TaskIndex:                 r.TaskRaw.TaskIndex,
		RequestId:                 r.TaskRaw.RequestID,
		FeeToken:                  r.TaskRaw.FeeToken,
		Payment:                   payment,
		RequestData:               r.TaskRaw.RequestData,
		CallbackAddress:           r.TaskRaw.CallbackAddress,
		CallbackFunctionId:        r.TaskRaw.CallbackFunctionID,
		TaskCreatedBlock:          r.TaskRaw.TaskCreatedBlock,
		QuorumNumbers:             r.TaskRaw.QuorumNumbers,
		QuorumThresholdPercentage: r.TaskRaw.QuorumThresholdPercentage,
	}
	price, ok := math.NewIntFromString(r.PriceFeedResponse.Price)
	if !ok {
		return nil, fmt.Errorf("failed to convert price: %s", r.PriceFeedResponse.Price)
	}

	priceFeedResponse := &pricetypes.PriceFeedResponse{
		ReferenceTaskIndex: r.PriceFeedResponse.ReferenceTaskIndex,
		Price:              price,
	}

	validatedData := &types.RequestPostRequestValidatedData{}
	if err := proto.Unmarshal(r.ValidatedData, validatedData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal validated data: %v", err)
	}

	return &pricetypes.MsgVoteFinalizedRequestPrice{
		TaskRaw:           taskRaw,
		ValidatedData:     validatedData,
		PriceFeedResponse: priceFeedResponse,
	}, nil
}

// FromProtoMessage converts MsgVoteFinalizedRequestPrice to RPCVoteFinalizedRequestPrice
func MsgVoteFinalizedRequestPriceFromProtoMessage(msg *pricetypes.MsgVoteFinalizedRequestPrice) (*RPCVoteFinalizedRequestPrice, error) {
	if msg == nil {
		return nil, fmt.Errorf("nil MsgVoteFinalizedRequestPrice")
	}

	validatedData, err := proto.Marshal(msg.ValidatedData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal validated data: %v", err)
	}

	return &RPCVoteFinalizedRequestPrice{
		TaskRaw: &RPCTaskRaw{
			TaskIndex:                 msg.TaskRaw.TaskIndex,
			RequestID:                 msg.TaskRaw.RequestId,
			FeeToken:                  msg.TaskRaw.FeeToken,
			Payment:                   msg.TaskRaw.Payment.String(),
			RequestData:               msg.TaskRaw.RequestData,
			CallbackAddress:           msg.TaskRaw.CallbackAddress,
			CallbackFunctionID:        msg.TaskRaw.CallbackFunctionId,
			TaskCreatedBlock:          msg.TaskRaw.TaskCreatedBlock,
			QuorumNumbers:             msg.TaskRaw.QuorumNumbers,
			QuorumThresholdPercentage: msg.TaskRaw.QuorumThresholdPercentage,
		},
		ValidatedData: validatedData,
		PriceFeedResponse: &RPCPriceFeedResponse{
			ReferenceTaskIndex: msg.PriceFeedResponse.ReferenceTaskIndex,
			Price:              msg.PriceFeedResponse.Price.String(),
		},
	}, nil
}
