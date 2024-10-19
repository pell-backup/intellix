package app

import (
	"context"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	pricetypes "intellix/x/price/dvs/types"

	"github.com/0xPellNetwork/pelldvs/aggregator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/crypto/sha3"

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
		Reponse: resData,
		// TODO: move calcDigest to biz
		ResponseDigest: calcDigest(resData),
	}, err
}

func convertValidatedResponse(validatedData *aggregator.ValidatedResponse) *pricetypes.RequestPostRequestPriceFeedValidatedData {
	var nonSignersPubkeysG1 []*pricetypes.G1Point
	for _, pubkey := range validatedData.NonSignersPubkeysG1 {
		x := pubkey.X.Bytes()
		y := pubkey.Y.Bytes()
		nonSignersPubkeysG1 = append(nonSignersPubkeysG1, &pricetypes.G1Point{
			X: x[:],
			Y: y[:],
		})
	}

	var quorumApksG1 []*pricetypes.G1Point
	for _, pubkey := range validatedData.QuorumApksG1 {
		x := pubkey.X.Bytes()
		y := pubkey.Y.Bytes()
		quorumApksG1 = append(quorumApksG1, &pricetypes.G1Point{
			X: x[:],
			Y: y[:],
		})
	}

	var signersApkG2 *pricetypes.G2Point
	if validatedData.SignersApkG2 != nil {
		xReal := validatedData.SignersApkG2.X.A0.Bytes()
		xImag := validatedData.SignersApkG2.X.A1.Bytes()
		yReal := validatedData.SignersApkG2.Y.A0.Bytes()
		yImag := validatedData.SignersApkG2.Y.A1.Bytes()
		signersApkG2 = &pricetypes.G2Point{
			XReal: xReal[:],
			XImag: xImag[:],
			YReal: yReal[:],
			YImag: yImag[:],
		}
	}

	var signersAggSigG1 *pricetypes.Signature
	if validatedData.SignersAggSigG1 != nil {
		s := validatedData.SignersAggSigG1.Bytes()
		signersAggSigG1 = &pricetypes.Signature{
			Sig: s[:],
		}
	}

	var nonSignerStakeIndices []*pricetypes.UInt32List
	for _, stakeIndices := range validatedData.NonSignerStakeIndices {
		nonSignerStakeIndices = append(nonSignerStakeIndices, &pricetypes.UInt32List{
			Values: stakeIndices,
		})
	}

	resp := &pricetypes.RequestPostRequestPriceFeedValidatedData{
		Data:                         validatedData.Data,
		Error:                        validatedData.Err.Error(),
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
	data, err := handlerSrc.InvokeRouterByData(sdkCtx, req.Request.Data, convertValidatedResponse(&req.Response))
	if err != nil {
		return nil, err
	}

	return &avsi.ResponsePostRequest{
		Receipt: data,
	}, nil
}
