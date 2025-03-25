package server

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"cosmossdk.io/math"
	"github.com/ontio/ontology-crypto/keypair"

	"intellix/common"
	"intellix/dvs/vrf/types"
)

// HandleVRFRandomNumberRequest handles the VRF random number request
// It verifies the VRF proof and returns the random number
// The random number is truncated to 255 bits to ensure it is within the defined range
// The random number is then returned in the response
func (s *Server) HandleVRFRandomNumberRequest(ctx context.Context, request *types.VRFTaskRequest) (*types.VRFTaskResponse, error) {
	s.logger.Info("HandleVRFRandomNumberRequest", "in", fmt.Sprintf("%+v", request))

	// Decode the public key from config
	pubKeyStr := strings.TrimPrefix(s.eccKeyPair.ECCPublicKey, "0x")
	pubKeyBuf, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		s.logger.Error("Failed to decode public key", "error", err)
		return nil, fmt.Errorf("Failed to decode public key")
	}

	// Deserialize the public key
	pubKey, err := keypair.DeserializePublicKey(pubKeyBuf)
	if err != nil {
		s.logger.Error("Failed to deserialize public key", "error", err)
		return nil, err
	}

	// Verify the VRF proof and return the random number
	var randomNumbers []math.Int
	for _, v := range request.VrfData {
		// Verify the VRF proof
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

	s.logger.Info("HandleVRFRandomNumberRequest", "randomNumbers", randomNumbers, "taskIndex", request.TaskMetadata.TaskIndex, "requestId", request.TaskMetadata.RequestId)

	return &types.VRFTaskResponse{
		TaskIndex:    request.TaskMetadata.TaskIndex,
		RequestId:    request.TaskMetadata.RequestId,
		RandomNumber: randomNumbers,
	}, nil
}

func bytesToNumber(b []byte) *big.Int {
	num := new(big.Int)
	num.SetBytes(b)
	return num
}
