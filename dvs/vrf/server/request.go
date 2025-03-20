package server

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ontio/ontology-crypto/keypair"
	"intellix/common"
	"intellix/dvs/vrf/types"
	"math/big"
	"strings"

	"cosmossdk.io/math"
)

func (s Server) HandleVRFRandomNumberRequest(ctx context.Context, request *types.VRFTaskRequest) (*types.VRFTaskResponse, error) {
	s.logger.Info("HandleVRFRandomNumberRequest", "in", fmt.Sprintf("%+v", request))

	// TODO: use config
	pubKeyStr := "03515b41d933f87408a1837681b62a1b49521e10e614027685eac41fe66d1f5c55"
	if strings.HasPrefix(pubKeyStr, "0x") {
		pubKeyStr = pubKeyStr[2:]
	}
	pubKeyBuf, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		s.logger.Error("Failed to decode private key", "error", err)
	}
	pubKey, err := keypair.DeserializePublicKey(pubKeyBuf)
	if err != nil {
		s.logger.Error("Failed to deserialize public key", "error", err)
		return nil, err
	}

	var randomNumbers []math.Int
	for _, v := range request.VrfData {
		verified, err := common.VerifyVRF(pubKey, request.TaskMetadata, v.VrfValue, v.VrfProof)
		if err != nil {
			s.logger.Error("Failed to verify VRF", "error", err)
			return nil, err
		}

		if !verified {
			s.logger.Error("VRF verification failed")
			return nil, err
		}

		randomNumber := bytesToNumber(v.VrfValue)

		// Define a 255-bit limit (2^255) as the maximum value for our random number
		// VRF output can potentially be larger than needed, so we truncate it
		limit := new(big.Int).Lsh(big.NewInt(1), 255) // 1 << 255
		// Apply modulo operation to ensure the number is within our defined range
		// This truncates the bigint to prevent overflow and ensure consistent bit length
		randomNumber.Mod(randomNumber, limit)

		value := math.NewIntFromBigInt(randomNumber)
		randomNumbers = append(randomNumbers, value)
	}

	s.logger.Info("HandleVRFRandomNumberRequest", "randomNumbers", randomNumbers)

	return &types.VRFTaskResponse{
		TaskIndex:    request.TaskMetadata.TaskIndex,
		RandomNumber: randomNumbers,
	}, nil
}

func bytesToNumber(b []byte) *big.Int {
	num := new(big.Int)
	num.SetBytes(b)
	return num
}
