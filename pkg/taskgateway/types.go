package taskgateway

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"intellix/pkg/pelldvs/types"
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

/*
func convertPbToBN254G1Point(pb *types.G1Point) *priceOracle.BN254G1Point {
	return &priceOracle.BN254G1Point{
		X: big.NewInt(0).SetBytes(pb.X),
		Y: big.NewInt(0).SetBytes(pb.Y),
	}
}

func convertPbToBN254G1PointList(pb []*types.G1Point) []priceOracle.BN254G1Point {
	list := make([]priceOracle.BN254G1Point, len(pb))
	for i, p := range pb {
		list[i] = *convertPbToBN254G1Point(p)
	}
	return list
}

func convertPbToBN254G2Point(pb *types.G2Point) *priceOracle.BN254G2Point {
	return &priceOracle.BN254G2Point{
		X: [2]*big.Int{
			big.NewInt(0).SetBytes(pb.XReal),
			big.NewInt(0).SetBytes(pb.XImag),
		},
		Y: [2]*big.Int{
			big.NewInt(0).SetBytes(pb.YReal),
			big.NewInt(0).SetBytes(pb.YImag),
		},
	}
}
*/

func convertUInt32ListToSlice(list []*types.NonSignerStakeIndice) [][]uint32 {
	slice := make([][]uint32, len(list))
	for i, l := range list {
		slice[i] = l.NonSignerStakeIndice
	}
	return slice
}

type RespondToTaskResponse struct {
	Error string `json:"error"`
}
