package pelldvs

import (
	"context"
	"github.com/cometbft/cometbft/rpc/client/local"
	"github.com/cosmos/cosmos-sdk/client"
)

type Client struct {
	ctx client.Context
}

func NewClient(clientCtx client.Context) (*Client, error) {
	return &Client{
		ctx: clientCtx,
	}, nil
}

func (c *Client) WithTendermintClient(tmNode *local.Local) *Client {
	c.ctx = c.ctx.WithClient(tmNode)
	return c
}

func (c *Client) RequestDVS(ctx context.Context, request []byte) error {
	// TODO: convert bytes to *avsiTypes.DVSRequest and call pell-dvs
	return nil
}
