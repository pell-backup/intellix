package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"intellix/x/processor/types"
)

func (k Keeper) ListProcessor(ctx context.Context, req *types.QueryListProcessorRequest) (*types.QueryListProcessorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProcessorKey))

	var processors []types.Processor
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var processor types.Processor
		if err := k.cdc.Unmarshal(value, &processor); err != nil {
			return err
		}

		processors = append(processors, processor)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListProcessorResponse{Processor: processors, Pagination: pageRes}, nil
}
