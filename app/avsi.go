package app

import (
	"context"
	"cosmossdk.io/math"
	"fmt"
	"github.com/cosmos/gogoproto/proto"
	pricetypes "intellix/x/price/dvs/types"

	avsi "github.com/0xPellNetwork/pelldvs/application"
)

func (app *App) ProcessRequest(ctx context.Context, req *avsi.RequestProcessRequest) (*avsi.ResponseProcessRequest, error) {
	var taskReq pricetypes.TaskRequest
	err := proto.Unmarshal(req.Request.Data, &taskReq)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal TaskRequest: %w", err)
	}

	// route
	response, err := app.DvsServer.ProcessDVSRequest(ctx, &pricetypes.RequestProcessDVSRequest{
		Task:    &taskReq,
		Height:  req.Request.Height,
		ChainId: math.NewIntFromBigInt(req.Request.ChainID),
	})
	if err != nil {
		return nil, fmt.Errorf("ProcessDVSRequest failed to process request: %w", err)
	}

	priceData, err := response.Price.Marshal()
	if err != nil {
		return nil, fmt.Errorf("ProcessDVSRequest price data marshal failed: %w", err)
	}

	var resp = &avsi.ResponseProcessRequest{
		Reponse:        priceData,
		ResponseDigest: response.PriceDigest,
	}

	return resp, nil
}

func (app *App) PostRequest(ctx context.Context, req *avsi.RequestPostRequest) (*avsi.ResponsePostRequest, error) {
	var taskReq pricetypes.TaskRequest
	err := proto.Unmarshal(req.Response.Data, &taskReq)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal TaskRequest: %w", err)
	}

	_, err = app.DvsServer.PostRequest(ctx, &pricetypes.RequestPostRequest{
		Task: &taskReq,
	})
	if err != nil {
		return nil, fmt.Errorf("PostRequest failed to process request: %w", err)
	}

	return &avsi.ResponsePostRequest{}, nil
}
