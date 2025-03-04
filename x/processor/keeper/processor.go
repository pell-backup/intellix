package keeper

import (
	"encoding/binary"
	"errors"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"intellix/x/processor/types"
)

func (k Keeper) AppendProcessor(ctx sdk.Context, processor types.Processor) (uint64, error) {
	if processor.ProcessorType != types.ProcessorType_WASM {
		return 0, errors.New("processor type must be WASM")
	}
	if processor.WasmCode == nil {
		return 0, errors.New("wasm code must be provided")
	}

	count := k.GetProcessorCount(ctx)
	processor.Id = count
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProcessorKey))
	value := k.cdc.MustMarshal(&processor)
	store.Set(GetProcessorIDBytes(processor.Id), value)
	k.SetProcessorCount(ctx, count+1)
	return count, nil
}

func (k Keeper) GetProcessorCount(ctx sdk.Context) uint64 {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.ProcessorCountKey)
	count := store.Get(byteKey)
	if count == nil {
		return 0
	}
	return binary.BigEndian.Uint64(count)
}

func GetProcessorIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

func (k Keeper) SetProcessorCount(ctx sdk.Context, count uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.ProcessorCountKey)
	store.Set(byteKey, GetProcessorIDBytes(count))
}

func (k Keeper) GetProcessor(ctx sdk.Context, id uint64) (val types.Processor, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProcessorKey))
	bz := store.Get(GetProcessorIDBytes(id))
	if bz == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(bz, &val)
	return val, true
}
