package taskgateway

import (
	"fmt"
	"net/rpc"

	"github.com/0xPellNetwork/pelldvs/libs/log"
)

// Client represents RPC client
type Client struct {
	client  *rpc.Client
	logger  log.Logger
	address string
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
	return &Client{
		client:  client,
		logger:  logger,
		address: address,
	}, nil
}

// reconnect attempts to reconnect to the RPC server
func (c *Client) reconnect() error {
	if c.client != nil {
		c.client.Close()
	}

	client, err := rpc.Dial("tcp", c.address)
	if err != nil {
		c.logger.Error("Failed to reconnect to RPC server", "error", err)
		return err
	}

	c.client = client
	c.logger.Info("Reconnected to RPC server", "address", c.address)
	return nil
}

func (c *Client) RespondToTask(req *RPCVoteFinalizedRequestIn) error {
	resp := &RespondToTaskResponse{}

	err := c.client.Call("TaskGateway.RespondToTask", req, resp)
	if err != nil {
		// retry
		if err.Error() == "connection is shut down" {
			c.logger.Info("Connection lost, attempting to reconnect...")
			if err := c.reconnect(); err != nil {
				return fmt.Errorf("failed to reconnect: %v", err)
			}
			err = c.client.Call("TaskGateway.RespondToTask", req, resp)
			if err != nil {
				c.logger.Error("RPC call failed after reconnection", "error", err.Error())
				return err
			}
		} else {
			c.logger.Error("RPC call failed", "error", err.Error())
			return err
		}
	}

	if resp.Error != "" {
		c.logger.Error("task RespondToTask failed", "error", resp.Error)
		return fmt.Errorf("task RespondToTask failed: %s", resp.Error)
	}

	c.logger.Info("Task response sent successfully", "TaskIndex", req.TaskRaw.TaskIndex)

	return nil
}
