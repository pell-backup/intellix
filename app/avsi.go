package app

import (
	"context"
	"github.com/0xPellNetwork/pelldvs/aggregator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/crypto/sha3"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
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
	// new SDK context
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx = sdkCtx.WithBlockHeight(req.Request.Height)
	sdkCtx = sdkCtx.WithChainID(req.Request.ChainID.String())

	handlerSrc := dvsservermanager.GetProcessRequestHandlerSrc()
	resData, err := handlerSrc.InvokeRouterByData(sdkCtx, req.Request.Data)
	return &avsi.ResponseProcessRequest{
		Reponse:        resData,
		ResponseDigest: nil,
	}, err
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
	// new SDK context
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx = sdkCtx.WithBlockHeight(req.Request.Height)
	sdkCtx = sdkCtx.WithChainID(req.Request.ChainID.String())

	handlerSrc := dvsservermanager.GetPostProcessRequestHandlerSrc()
	data, err := handlerSrc.InvokeRouterByData(sdkCtx, req.Request.Data, &pricetypes.RequestPostRequestPriceFeedValidatedData{
		// TODO: fill data
	})
	if err != nil {
		return nil, err
	}

	return &avsi.ResponsePostRequest{
		Receipt: data,
	}, nil
}
