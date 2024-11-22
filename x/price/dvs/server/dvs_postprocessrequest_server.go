package server

import (
	"context"
	"cosmossdk.io/math"
	"fmt"
	contractDataOracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/ethereum/go-ethereum/accounts/abi"
	pkgcontext "intellix/pkg/context"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	dvstypes "intellix/pkg/pelldvs/types"
	"intellix/x/price/dvs/types"
	pricetypes "intellix/x/price/types"
	"math/big"
)

type DvsPostProcessRequestServer struct {
	Server
}

func NewDvsPostProcessRequestServer(server Server) types.DvsPostProcessRequestServer {
	return &DvsPostProcessRequestServer{Server: server}
}

func (d DvsPostProcessRequestServer) decodePackedPriceFeedData(data []byte) (*contractDataOracle.IDataOracleTaskResponse, error) {
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

	r, ok := values[0].(struct {
		ReferenceTaskIndex uint32   `json:"referenceTaskIndex"`
		Price              *big.Int `json:"price"`
	})
	if !ok {
		return nil, fmt.Errorf("expected %T, got %T", &contractDataOracle.IDataOracleTaskResponse{}, values[0])
	}

	return &contractDataOracle.IDataOracleTaskResponse{
		ReferenceTaskIndex: r.ReferenceTaskIndex,
		Price:              r.Price,
	}, nil
}

func (d DvsPostProcessRequestServer) getDvsRequestValidatedData(ctx pkgcontext.Context) (*dvstypes.RequestPostRequestValidatedData, error) {
	reqData, ok := ctx.DvsPostResponseData()
	if !ok {
		return nil, fmt.Errorf("not DvsRequestData found")
	}
	validatedDataMsg, err := dvsservermanager.DecodeMsg(reqData)
	if err != nil {
		return nil, err
	}
	validatedData, ok := validatedDataMsg.(*dvstypes.RequestPostRequestValidatedData)
	if !ok {
		return nil, fmt.Errorf("expected %T, got %T", &dvstypes.RequestPostRequestValidatedData{}, validatedData)
	}
	return validatedData, nil
}

func (d DvsPostProcessRequestServer) PostProcessRequestPriceFeed(ctx context.Context, in *types.ProcessRequestPriceFeedIn) (*types.PostProcessRequestPriceFeedOut, error) {
	pkgCtx := pkgcontext.UnwrapContext(ctx)

	validatedData, err := d.getDvsRequestValidatedData(pkgCtx)
	if err != nil {
		return nil, err
	}
	taskResp, err := d.decodePackedPriceFeedData(validatedData.Data)
	if err != nil {
		return nil, err
	}

	// send VoteFinalizedRequestPrice Tx
	msg, err := d.sendVoteFinalizedRequestPriceTx(pkgCtx, in, validatedData, taskResp)
	if err != nil {
		return nil, err
	}

	// send gateway
	err = d.sendResponseToGateway(msg)
	if err != nil {
		return nil, err
	}

	return &types.PostProcessRequestPriceFeedOut{}, nil
}

func (d DvsPostProcessRequestServer) sendVoteFinalizedRequestPriceTx(ctx pkgcontext.Context, raw *types.ProcessRequestPriceFeedIn, validatedData *dvstypes.RequestPostRequestValidatedData, priceData *contractDataOracle.IDataOracleTaskResponse) (*pricetypes.MsgVoteFinalizedRequestPrice, error) {
	msg := &pricetypes.MsgVoteFinalizedRequestPrice{
		TaskRaw: &pricetypes.TaskRaw{
			TaskIndex:                 raw.Task.TaskIndex,
			RequestId:                 raw.Task.RequestId,
			FeeToken:                  raw.Task.FeeToken,
			Payment:                   raw.Task.Payment,
			RequestData:               raw.Task.RequestData,
			CallbackAddress:           raw.Task.CallbackAddress,
			CallbackFunctionId:        raw.Task.CallbackFunctionId,
			TaskCreatedBlock:          raw.Task.TaskCreatedBlock,
			QuorumNumbers:             raw.Task.QuorumNumbers,
			QuorumThresholdPercentage: raw.Task.QuorumThresholdPercentage,
		},
		ValidatedData: validatedData,
		PriceFeedResponse: &pricetypes.PriceFeedResponse{
			ReferenceTaskIndex: priceData.ReferenceTaskIndex,
			Price:              math.NewIntFromBigInt(priceData.Price),
		},
	}

	if err := d.Server.SignAndBroadcastTx(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (d DvsPostProcessRequestServer) sendResponseToGateway(price *pricetypes.MsgVoteFinalizedRequestPrice) error {
	return d.Server.taskGatewayClient.RespondToTask(price)
}
