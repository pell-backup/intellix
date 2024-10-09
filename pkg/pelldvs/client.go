package pelldvs

import (
	"context"
	"github.com/cometbft/cometbft/rpc/client/local"
	"github.com/cosmos/cosmos-sdk/client"
)

type Client struct {
	ctx client.Context
}

func NewClient(clientCtx client.Context) *Client {
	return &Client{
		ctx: clientCtx,
	}
}

func (c *Client) WithTendermintClient(tmNode *local.Local) *Client {
	c.ctx = c.ctx.WithClient(tmNode)
	return c
}

func (c *Client) RequestDVS(ctx context.Context, request []byte) error {
	// todo: convert bytes to *avsiTypes.DVSRequest and call pell-dvs
	return nil
}
