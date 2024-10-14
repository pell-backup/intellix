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
	var (
		data       []byte
		dataDigest []byte
	)

	switch taskReq.TaskType {
	case pricetypes.TaskType_PRICE_FEED:
		response, err := app.DvsServer.ProcessRequestPriceFeed(ctx, &pricetypes.RequestProcessRequestPriceFeed{
			Task:    &taskReq,
			Height:  req.Request.Height,
			ChainId: math.NewIntFromBigInt(req.Request.ChainID),
		})
		if err != nil {
			return nil, fmt.Errorf("ProcessDVSRequest failed to process request: %w", err)
		}
		data, err = response.Price.Marshal()
		if err != nil {
			return nil, fmt.Errorf("ProcessDVSRequest price data marshal failed: %w", err)
		}
		dataDigest = response.PriceDigest
	default:
		return nil, fmt.Errorf("unknown task type: %d", taskReq.TaskType)
	}

	var resp = &avsi.ResponseProcessRequest{
		Reponse:        data,
		ResponseDigest: dataDigest,
	}

	return resp, nil
}

func (app *App) PostRequest(ctx context.Context, req *avsi.RequestPostRequest) (*avsi.ResponsePostRequest, error) {
	var taskReq pricetypes.TaskRequest
	err := proto.Unmarshal(req.Response.Data, &taskReq)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal TaskRequest: %w", err)
	}

	switch taskReq.TaskType {
	case pricetypes.TaskType_PRICE_FEED:
		_, err = app.DvsServer.PostRequestPriceFeed(ctx, &pricetypes.RequestPostRequestPriceFeed{
			Task: &taskReq,
		})
		if err != nil {
			return nil, fmt.Errorf("PostRequest failed to process request: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown task type: %d", taskReq.TaskType)
	}

	return &avsi.ResponsePostRequest{}, nil
}
