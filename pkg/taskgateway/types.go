package taskgateway

import (
	"fmt"
	"github.com/0xPellNetwork/pelldvs/crypto/bls"
	dataOracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/ethereum/go-ethereum/common"
	"intellix/pkg/pelldvs/types"
	"math/big"
)

type TaskGatewayCfg struct {
	SenderAddress    string `mapstructure:"sender_address"`
	EthEndpoint      string `mapstructure:"eth_endpoint"`
	BftNetworkRemote string `mapstructure:"bft_network_remote"`
	ContractAddress  string `mapstructure:"contract_address"`
}

func (t TaskGatewayCfg) Validate() error {
	if t.EthEndpoint == "" {
		return fmt.Errorf("eth endpoint cannot be empty")
	}
	if t.BftNetworkRemote == "" {
		return fmt.Errorf("cometbft network remote cannot be empty")
	}
	if t.ContractAddress == "" {
		return fmt.Errorf("contract address cannot be empty")
	}
	if t.SenderAddress == "" {
		return fmt.Errorf("sender address cannot be empty")
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
