package server

import (
	"context"
	"cosmossdk.io/math"
	"fmt"
	contractPriceOracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/PriceOracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	dvstypes "intellix/pkg/pelldvs/types"
	"intellix/x/price/dvs/types"
	"time"
)

type DvsPostProcessRequestServer struct {
	Server
}

func NewDvsPostProcessRequestServer(server Server) types.DvsPostProcessRequestServer {
	return &DvsPostProcessRequestServer{Server: server}
}

func (d DvsPostProcessRequestServer) decodePackedPriceFeedData(data []byte) (*contractPriceOracle.IPriceOracleTaskResponse, error) {
	taskResponseType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "referenceTaskIndex",
			Type: "uint32",
		},
		{
			Name: "price",
			Type: "uint256",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ABI type: %w", err)
	}

	arguments := abi.Arguments{
		{
			Type: taskResponseType,
		},
	}

	// decode
	values, err := arguments.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack data: %w", err)
	}

	if len(values) != 1 {
		return nil, fmt.Errorf("unexpected number of values: got %d, want 1", len(values))
	}

	r, ok := values[0].(contractPriceOracle.IPriceOracleTaskResponse)
	if !ok {
		return nil, fmt.Errorf("expected %T, got %T", &contractPriceOracle.IPriceOracleTaskResponse{}, r)
	}

	return &r, nil
}

func (d DvsPostProcessRequestServer) PostProcessRequestPriceFeed(ctx context.Context, msg *dvstypes.RequestPostRequestValidatedData) (*types.ResponsePostRequestPriceFeed, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	reqData, ok := dvsservermanager.CtxGetDvsRequestData(sdkCtx)
	if !ok {
		return nil, fmt.Errorf("not DvsRequestData found")
	}
	requestMsg, err := dvsservermanager.DecodeMsg(reqData)
	if err != nil {
		return nil, err
	}
	taskRequest, ok := requestMsg.(*types.ProcessPriceFeedMsg)
	if !ok {
		return nil, fmt.Errorf("expected %T, got %T", &types.ProcessPriceFeedMsg{}, taskRequest)
	}

	taskResp, err := d.decodePackedPriceFeedData(msg.Data)
	if err != nil {
		return nil, err
	}

	// just send VoteFinalizedRequestPrice Tx
	err = d.sendVoteFinalizedRequestPriceTx(sdkCtx, taskRequest.Raw, taskResp)
	if err != nil {
		return nil, err
	}

	return &types.ResponsePostRequestPriceFeed{}, nil
}

func (d DvsPostProcessRequestServer) sendVoteFinalizedRequestPriceTx(ctx sdk.Context, raw *types.TaskRequestRaw, priceData *contractPriceOracle.IPriceOracleTaskResponse) error {
	msg := &types.MsgVoteFinalizedRequestPrice{
		TaskIndex:                 priceData.ReferenceTaskIndex,
		RequestId:                 raw.RequestId,
		Price:                     math.LegacyNewDecFromBigInt(priceData.Price),
		Timestamp:                 time.Now().Unix(),
		FeeToken:                  raw.FeeToken,
		Payment:                   raw.Payment,
		CallbackAddress:           raw.CallbackAddress,
		CallbackFunctionId:        raw.CallbackFunctionId,
		QuorumNumbers:             raw.QuorumNumbers,
		QuorumThresholdPercentage: raw.QuorumThresholdPercentage,
	}

	if err := d.Server.SignAndBroadcastTx(ctx, msg); err != nil {
		return err
	}
	return nil
}
