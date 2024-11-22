package taskgateway

import (
	"fmt"
	"intellix/x/price/types"
	"net/rpc"

	"github.com/0xPellNetwork/pelldvs/libs/log"
)

// Client represents RPC client
type Client struct {
	client *rpc.Client
	logger log.Logger
}

// NewClient creates a new TaskGateway RPC client
func NewClient(address string, logger log.Logger) (*Client, error) {
	if logger == nil {
		logger = log.NewNopLogger()
	}

	client, err := rpc.Dial("tcp", address)
	if err != nil {
		logger.Error("Failed to connect to RPC server", "error", err)
		return nil, fmt.Errorf("failed to connect to RPC server: %v", err)
	}

	logger.Info("Connected to RPC server", "address", address)
	return &Client{client: client}, nil
}

func (c *Client) RespondToTask(req *types.MsgVoteFinalizedRequestPrice) error {
	resp := &RespondToTaskResponse{}

	err := c.client.Call("Server.RespondToTask", req, resp)
	if err != nil {
		c.logger.Error("RPC call failed", "error", err)
		return err
	}

	if resp.Error != "" {
		c.logger.Error("task RespondToTask failed", "error", resp.Error)
		return fmt.Errorf("task RespondToTask failed: %s", resp.Error)
	}

	c.logger.Info("Task response sent successfully", "TaskIndex", req.TaskRaw.TaskIndex)

	return nil
}
