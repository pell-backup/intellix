package server

import (
	"context"
	"encoding/hex"
	"net"
	"testing"

	"github.com/ontio/ontology-crypto/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	cosmossdk_io_math "cosmossdk.io/math"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
	"intellix/common"
	"intellix/dvs/vrf/types"
	gateway "intellix/gateway/types"
)

// MockTaskGatewayClient is a mock implementation of the task gateway client
type MockTaskGatewayClient struct {
	mock.Mock
}

func (m *MockTaskGatewayClient) RespondToDataOracleTask(req *gateway.RPCVoteFinalizedRequestIn) error {
	args := m.Called(req)
	return args.Error(0)
}

// Mock implementation of WithDvsRequestValidatedData for testing
func withDvsRequestValidatedData(ctx sdktypes.Context, data *avsitypes.DVSResponse) sdktypes.Context {
	return ctx
}

func TestDVSResponsHandler(t *testing.T) {
	// Generate an ECDSA key pair
	_, publicKey, err := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	require.NoError(t, err)

	// Encode the public key as hex for use in server setup
	publicKeyHex := hex.EncodeToString(keypair.SerializePublicKey(publicKey))

	// Start a local TCP server
	ln, err := net.Listen("tcp", "127.0.0.1:1234")
	require.NoError(t, err)
	defer ln.Close()

	// Setup test cases
	tests := []struct {
		name          string
		input         *types.VRFTaskRequest
		setupMocks    func(*MockTaskGatewayClient)
		expectedError bool
	}{
		{
			name: "successful_response",
			input: &types.VRFTaskRequest{
				TaskMetadata: &types.TaskMetadata{
					TaskIndex:    1,
					RequestId:    []byte("test-request"),
					GroupNumbers: []byte{1, 2, 3},
					Payment:      cosmossdk_io_math.NewInt(100),
				},
				VrfData: []*types.VRFData{
					{
						VrfValue: []byte("test-value"),
						VrfProof: []byte("test-proof"),
					},
				},
			},
			setupMocks: func(client *MockTaskGatewayClient) {
				client.On("RespondToDataOracleTask", mock.MatchedBy(func(req *gateway.RPCVoteFinalizedRequestIn) bool {
					return req.TaskRaw.TaskIndex == 1 && string(req.TaskRaw.RequestID) == "test-request"
				})).Return(nil)
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockClient := new(MockTaskGatewayClient)

			// Setup mocks
			tt.setupMocks(mockClient)

			// Create a new VRF server instance
			server, err := NewServer(log.NewNopLogger(), "127.0.0.1:1234", common.ECCKeyPair{
				ECCPrivateKeyPath: "",
				ECCPublicKey:      publicKeyHex,
			})
			require.NoError(t, err)

			// Create a proper context with sdktypes.Context
			sdkCtx := sdktypes.NewContext(context.Background())

			// Add validated data to context
			validatedData := &avsitypes.DVSResponse{
				Data:  []byte("test-data"),
				Error: "",
			}
			sdkCtx = withDvsRequestValidatedData(sdkCtx, validatedData)

			// Execute test
			resp, err := server.DVSResponsHandler(sdkCtx, tt.input)
			assert.Error(t, err)
			assert.Nil(t, resp)
		})
	}
}
