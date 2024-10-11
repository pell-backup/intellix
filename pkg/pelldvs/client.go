package pelldvs

import (
	"context"

	"github.com/cometbft/cometbft/rpc/client/local"
	"github.com/cosmos/cosmos-sdk/client"

	avsi "github.com/0xPellNetwork/pelldvs/application"
)

type Client struct {
	ctx client.Context
}

// TODO: create a new client from the local client in pelldvs
func NewClient(clientCtx client.Context) (*Client, error) {
	return &Client{
		ctx: clientCtx,
	}, nil
}

func (c *Client) WithTendermintClient(tmNode *local.Local) *Client {
	c.ctx = c.ctx.WithClient(tmNode)
	return c
}

// TODO: use the RequesstDVS interface from pelldvs localclient directly. No need to define it here.
func (c *Client) RequestDVS(ctx context.Context, request *avsi.RequestProcessRequest) error {
	return nil
}
