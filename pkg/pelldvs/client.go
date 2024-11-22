package pelldvs

import (
	"context"
	"fmt"
	"github.com/0xPellNetwork/pelldvs/avsi/types"
	"github.com/0xPellNetwork/pelldvs/libs/log"
	ctypes "github.com/0xPellNetwork/pelldvs/rpc/core/types"
	"github.com/0xPellNetwork/pelldvs/rpc/jsonrpc/client"

	"github.com/0xPellNetwork/pelldvs/config"
)

type Client struct {
	logger        log.Logger
	pellDVSClient *client.Client
}

// TODO: create a new client from the local client in pelldvs
func NewClient(
	logger log.Logger,
	remoteAddr string,
) (*Client, error) {
	var c = &Client{
		logger: logger,
	}
	if remoteAddr == "" {
		remoteAddr = config.DefaultConfig().RPC.ListenAddress
	}

	dvsClient, err := client.New(remoteAddr)
	if err != nil {
		c.logger.Error("Failed to create DVS client", "error", err)
		return nil, err
	}
	c.pellDVSClient = dvsClient
	return c, nil
}

// TODO: use the RequesstDVS interface from pelldvs localclient directly. No need to define it here.
func (c *Client) RequestDVS(ctx context.Context, request *types.RequestProcessRequest) error {
	if c.pellDVSClient == nil {
		return fmt.Errorf("pelldvs client has not been initialized")
	}
	result := ctypes.ResultDvsTask{}

	_, err := c.pellDVSClient.Call(ctx, "send_task", map[string]interface{}{
		"task":                         request.Request.Data,
		"height":                       request.Request.Height,
		"chainid":                      request.Request.ChainId,
		"quorum_numbers":               request.Request.QuorumNumbers,
		"quorum_threshold_percentages": request.Request.QuorumThresholdPercentages,
	}, &result)
	if err != nil {
		c.logger.Error("Failed to send dvs task", "error", err)
		return fmt.Errorf("failed to send dvs task: %v", err)
	}

	return nil
}
