package server

import (
	"context"
	"encoding/hex"
	"net"
	"testing"

	"github.com/ontio/ontology-crypto/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cosmossdk_io_math "cosmossdk.io/math"

	"github.com/0xPellNetwork/pelldvs-libs/log"
	"intellix/common"
	"intellix/dvs/vrf/types"
)

const (
	testAddr = "127.0.0.1:1234"
)

// TestBytesToNumber_ExtendedCases tests the conversion of various byte slices to numeric values.
func TestBytesToNumber_ExtendedCases(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "Empty input",
			input:    []byte{},
			expected: "0",
		},
		{
			name:     "Single byte",
			input:    []byte{0x01},
			expected: "1",
		},
		{
			name:     "Multiple bytes",
			input:    []byte{0x01, 0x02, 0x03},
			expected: "66051",
		},
		{
			name:     "High byte values",
			input:    []byte{0xFF, 0xFF},
			expected: "65535",
		},
		{
			name:     "32-byte array of zeros",
			input:    make([]byte, 32),
			expected: "0",
		},
		{
			name:     "Big endian example",
			input:    []byte{0x00, 0x01},
			expected: "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bytesToNumber(tt.input)
			assert.Equal(t, tt.expected, result.String())
		})
	}
}

// TestHandleVRFRandomNumberRequest tests both successful and failing scenarios of VRF request handling.
func TestHandleVRFRandomNumberRequest(t *testing.T) {
	// Generate an ECDSA key pair
	privateKey, publicKey, err := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	require.NoError(t, err)

	// Encode the public key as hex for use in server setup
	publicKeyHex := hex.EncodeToString(keypair.SerializePublicKey(publicKey))

	// Start a local TCP server
	ln, err := net.Listen("tcp", testAddr)
	require.NoError(t, err)
	defer ln.Close()

	// Create a new VRF server instance
	server, err := NewServer(log.NewNopLogger(), testAddr, common.ECCKeyPair{
		ECCPrivateKeyPath: "",
		ECCPublicKey:      publicKeyHex,
	})
	require.NoError(t, err)

	// Base metadata used in all VRF task requests
	baseTaskData := &types.TaskMetadata{
		TaskIndex:                0,
		RequestId:                nil,
		FeeToken:                 "",
		Payment:                  cosmossdk_io_math.Int{},
		RequestData:              nil,
		CallbackAddress:          "",
		CallbackFunctionId:       nil,
		TaskCreatedBlock:         0,
		GroupNumbers:             nil,
		GroupThresholdPercentage: 0,
		AdvanceDecode:            false,
	}

	// Compute correct VRF value and proof
	validValue, validProof, err := common.ComputeVRF(privateKey, baseTaskData)
	require.NoError(t, err)

	t.Run("Successful VRF verification", func(t *testing.T) {
		request := &types.VRFTaskRequest{
			TaskMetadata: baseTaskData,
			VrfData: []*types.VRFData{
				{
					VrfValue: validValue,
					VrfProof: validProof,
				},
			},
		}

		response, err := server.HandleVRFRandomNumberRequest(context.Background(), request)

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, request.TaskMetadata.RequestId, response.RequestId)
		require.Len(t, response.RandomNumber, 1)

		randomNum := response.RandomNumber[0].BigInt()
		require.NotZero(t, randomNum.BitLen())
		require.LessOrEqual(t, randomNum.BitLen(), 255)
	})

	t.Run("Nil request", func(t *testing.T) {
		resp, err := server.HandleVRFRandomNumberRequest(context.Background(), nil)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Empty VRF data array", func(t *testing.T) {
		resp, err := server.HandleVRFRandomNumberRequest(context.Background(), &types.VRFTaskRequest{
			TaskMetadata: baseTaskData,
			VrfData:      []*types.VRFData{},
		})
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Missing TaskMetadata", func(t *testing.T) {
		resp, err := server.HandleVRFRandomNumberRequest(context.Background(), &types.VRFTaskRequest{
			TaskMetadata: nil,
			VrfData: []*types.VRFData{
				{VrfValue: validValue, VrfProof: validProof},
			},
		})
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Invalid VRF proof", func(t *testing.T) {
		request := &types.VRFTaskRequest{
			TaskMetadata: baseTaskData,
			VrfData: []*types.VRFData{
				{
					VrfValue: validValue,
					VrfProof: []byte("invalid-proof-bytes"),
				},
			},
		}
		resp, err := server.HandleVRFRandomNumberRequest(context.Background(), request)
		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Tampered VRF value", func(t *testing.T) {
		request := &types.VRFTaskRequest{
			TaskMetadata: baseTaskData,
			VrfData: []*types.VRFData{
				{
					VrfValue: []byte("wrong-value"),
					VrfProof: validProof,
				},
			},
		}
		resp, err := server.HandleVRFRandomNumberRequest(context.Background(), request)
		require.Error(t, err)
		assert.Nil(t, resp)
	})
}
