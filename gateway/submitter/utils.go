package submitter

import (
	"fmt"
	"github.com/0xPellNetwork/pelldvs/crypto/bls"
	contractdataoracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/ethereum/go-ethereum/common"
	"intellix/gateway/types"
	"math/big"
)

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

func validateBLSComponents(data *types.RPCValidatedData) error {
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
