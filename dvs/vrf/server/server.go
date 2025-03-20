package server

import (
	"fmt"
	"intellix/common"
	"intellix/dvs/vrf/types"
	taskgateway "intellix/gateway/submitter"

	"github.com/0xPellNetwork/pelldvs-libs/log"
)

// Server struct represents the server with a logger and a chain connector client.
type Server struct {
	logger            log.Logger // Logger for logging messages.
	taskGatewayClient *taskgateway.Client
	eccKeyPair        common.ECCKeyPair
}

// NewServer creates a new Server instance with the provided logger and gateway RPC client URL.
func NewServer(
	logger log.Logger,
	gatewayAddr string,
	eccKeyPair common.ECCKeyPair,
) (Server, error) {

	s := Server{
		logger:     logger,
		eccKeyPair: eccKeyPair,
	}

	taskGatewayClient, err := taskgateway.NewClient(gatewayAddr, logger)
	if err != nil {
		return s, err
	}

	s.taskGatewayClient = taskGatewayClient
	return s, nil
}

// Logger returns a module-specific logger.
func (k *Server) Logger() log.Logger {
	// Add module-specific information to the logger.
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
