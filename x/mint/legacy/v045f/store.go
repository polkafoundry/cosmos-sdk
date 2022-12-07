package v045f

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
)

func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	pool := types.DefaultInitialPool()
	b, err := cdc.Marshal(&pool)
	if err != nil {
		return err
	}
	store.Set(types.PoolKey, b)
	return nil
}
