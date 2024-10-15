package app

import (
	"context"
	"fmt"
	"github.com/0xPellNetwork/pelldvs/aggregator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"golang.org/x/crypto/sha3"
	pricetypes "intellix/x/price/dvs/types"

	avsi "github.com/0xPellNetwork/pelldvs/application"
)

func calcDigest(data []byte) []byte {
	var taskResponseDigest [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(data)
	copy(taskResponseDigest[:], hasher.Sum(nil)[:32])

	return taskResponseDigest[:]
}

func (app *App) ProcessRequest(ctx context.Context, req *avsi.RequestProcessRequest) (*avsi.ResponseProcessRequest, error) {
	var taskReq pricetypes.ProcessPriceFeedMsg
	err := proto.Unmarshal(req.Request.Data, &taskReq)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal TaskRequest: %w", err)
	}

	// new SDK context
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx = sdkCtx.WithBlockHeight(req.Request.Height)
	sdkCtx = sdkCtx.WithChainID(req.Request.ChainID.String())

	handler := app.MsgServiceRouter().Handler(&taskReq)
	res, err := handler(sdkCtx, &taskReq)
	if err != nil {
		return nil, err
	}

	// TODO: check and convert proto
	//for _, resMsg := range res.MsgResponses {
	//	switch resMsg.TypeUrl {
	//		// check and convert proto
	//	}
	//}

	var resp = &avsi.ResponseProcessRequest{
		Reponse:        res.Data,
		ResponseDigest: calcDigest(res.Data),
	}

	return resp, nil
}

func convertValidatedResponse(validatedData *aggregator.ValidatedResponse) *pricetypes.RequestPostRequestPriceFeedValidatedData {
	resp := &pricetypes.RequestPostRequestPriceFeedValidatedData{}
	resp.Data = validatedData.Data
	resp.Error = validatedData.Err.Error()
	resp.Hash = validatedData.Hash
	resp.NonSignerQuorumBitmapIndices = validatedData.NonSignerQuorumBitmapIndices
	resp.QuorumApkIndices = validatedData.QuorumApkIndices
	resp.TotalStakeIndices = validatedData.TotalStakeIndices
	var nonSignerStakeIndices []*pricetypes.UInt32List
	for _, stakeIndices := range validatedData.NonSignerStakeIndices {
		nonSignerStakeIndices = append(nonSignerStakeIndices, &pricetypes.UInt32List{
			Values: stakeIndices,
		})
	}
	resp.NonSignerStakeIndices = nonSignerStakeIndices

	return resp
}

func (app *App) PostRequest(ctx context.Context, req *avsi.RequestPostRequest) (*avsi.ResponsePostRequest, error) {
	var taskReq pricetypes.PostPriceFeedMsg
	err := proto.Unmarshal(req.Request.Data, &taskReq)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal TaskRequest: %w", err)
	}
	// new SDK context
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx = sdkCtx.WithBlockHeight(req.Request.Height)
	sdkCtx = sdkCtx.WithChainID(req.Request.ChainID.String())

	// fill validated data
	taskReq.ValidatedData = convertValidatedResponse(&req.Response)

	handler := app.MsgServiceRouter().Handler(&taskReq)
	res, err := handler(sdkCtx, &taskReq)
	if err != nil {
		return nil, err
	}
	// TODO: check and convert proto

	return &avsi.ResponsePostRequest{
		Receipt: res.Data,
	}, nil
}
