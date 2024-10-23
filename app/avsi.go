package app

import (
	"context"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	dvstypes "intellix/pkg/pelldvs/types"

	"github.com/0xPellNetwork/pelldvs/aggregator"
	sdk "github.com/cosmos/cosmos-sdk/types"

	avsi "github.com/0xPellNetwork/pelldvs/application"
)

func (app *App) ProcessRequest(ctx context.Context, req *avsi.RequestProcessRequest) (*avsi.ResponseProcessRequest, error) {
	// new SDK context
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx = sdkCtx.WithBlockHeight(req.Request.Height)
	sdkCtx = sdkCtx.WithChainID(req.Request.ChainID.String())

	handlerSrc := dvsservermanager.GetProcessRequestHandlerSrc()
	res, err := handlerSrc.InvokeRouterRawByData(sdkCtx, req.Request.Data)
	if err != nil {
		return nil, err
	}

	return &avsi.ResponseProcessRequest{
		Reponse:        res.CustomData,
		ResponseDigest: res.CustomDigest,
	}, err
}

func convertValidatedResponse(validatedData *aggregator.ValidatedResponse) *dvstypes.RequestPostRequestValidatedData {
	var nonSignersPubkeysG1 []*dvstypes.G1Point
	for _, pubkey := range validatedData.NonSignersPubkeysG1 {
		x := pubkey.X.Bytes()
		y := pubkey.Y.Bytes()
		nonSignersPubkeysG1 = append(nonSignersPubkeysG1, &dvstypes.G1Point{
			X: x[:],
			Y: y[:],
		})
	}

	var quorumApksG1 []*dvstypes.G1Point
	for _, pubkey := range validatedData.QuorumApksG1 {
		x := pubkey.X.Bytes()
		y := pubkey.Y.Bytes()
		quorumApksG1 = append(quorumApksG1, &dvstypes.G1Point{
			X: x[:],
			Y: y[:],
		})
	}

	var signersApkG2 *dvstypes.G2Point
	if validatedData.SignersApkG2 != nil {
		xReal := validatedData.SignersApkG2.X.A0.Bytes()
		xImag := validatedData.SignersApkG2.X.A1.Bytes()
		yReal := validatedData.SignersApkG2.Y.A0.Bytes()
		yImag := validatedData.SignersApkG2.Y.A1.Bytes()
		signersApkG2 = &dvstypes.G2Point{
			XReal: xReal[:],
			XImag: xImag[:],
			YReal: yReal[:],
			YImag: yImag[:],
		}
	}

	var signersAggSigG1 *dvstypes.G1Point
	if validatedData.SignersAggSigG1 != nil {
		x := validatedData.SignersAggSigG1.X.Bytes()
		y := validatedData.SignersAggSigG1.Y.Bytes()
		signersAggSigG1 = &dvstypes.G1Point{
			X: x[:],
			Y: y[:],
		}
	}

	var nonSignerStakeIndices []*dvstypes.UInt32List
	for _, stakeIndices := range validatedData.NonSignerStakeIndices {
		nonSignerStakeIndices = append(nonSignerStakeIndices, &dvstypes.UInt32List{
			Values: stakeIndices,
		})
	}
	var errMsg string
	if validatedData.Err != nil {
		errMsg = validatedData.Err.Error()
	}

	resp := &dvstypes.RequestPostRequestValidatedData{
		Data:                         validatedData.Data,
		Error:                        errMsg,
		Hash:                         validatedData.Hash,
		NonSignersPubkeysG1:          nonSignersPubkeysG1,
		QuorumApksG1:                 quorumApksG1,
		SignersApkG2:                 signersApkG2,
		SignersAggSigG1:              signersAggSigG1,
		NonSignerQuorumBitmapIndices: validatedData.NonSignerQuorumBitmapIndices,
		QuorumApkIndices:             validatedData.QuorumApkIndices,
		TotalStakeIndices:            validatedData.TotalStakeIndices,
		NonSignerStakeIndices:        nonSignerStakeIndices,
	}

	return resp
}

func (app *App) PostRequest(ctx context.Context, req *avsi.RequestPostRequest) (*avsi.ResponsePostRequest, error) {
	// new SDK context
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx = sdkCtx.WithBlockHeight(req.Request.Height)
	sdkCtx = sdkCtx.WithChainID(req.Request.ChainID.String())

	handlerSrc := dvsservermanager.GetPostProcessRequestHandlerSrc()
	res, err := handlerSrc.InvokeRouterRawByData(sdkCtx, req.Request.Data, convertValidatedResponse(&req.Response))
	if err != nil {
		return nil, err
	}

	return &avsi.ResponsePostRequest{
		Receipt: res.CustomData,
	}, nil
}
