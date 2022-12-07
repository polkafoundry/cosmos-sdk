package mint

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint/keeper"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// fetch stored params
	params := k.GetParams(ctx)

	if !params.UsingFund {
		// fetch stored minter
		minter := k.GetMinter(ctx)
		// recalculate inflation rate
		totalStakingSupply := k.StakingTokenSupply(ctx)
		bondedRatio := k.BondedRatio(ctx)
		minter.Inflation = minter.NextInflationRate(params, bondedRatio)
		minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalStakingSupply)
		k.SetMinter(ctx, minter)
		// mint coins, update supply
		mintedCoin := minter.BlockProvision(params)
		mintedCoins := sdk.NewCoins(mintedCoin)

		err := k.MintCoins(ctx, mintedCoins)
		if err != nil {
			panic(err)
		}

		// send the minted coins to the fee collector account
		err = k.AddCollectedFees(ctx, mintedCoins)
		if err != nil {
			panic(err)
		}

		if mintedCoin.Amount.IsInt64() {
			defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCoin.Amount.Int64()), "minted_tokens")
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeMint,
				sdk.NewAttribute(types.AttributeKeyBondedRatio, bondedRatio.String()),
				sdk.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
				sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
			),
		)
	} else {
		// calculate the reward coin for a block
		blockReward := params.ValidatorRewardFund.QuoInt64(int64(params.FundDuration))

		// if the reward exceeds the remaining fund, use the remaining instead.
		pool := k.GetPool(ctx)
		remain := params.ValidatorRewardFund.Sub(pool.Used.ToDec())
		if remain.LT(blockReward) {
			blockReward = remain
		}

		blockRewardInt := blockReward.TruncateInt()
		pool.Used.Add(blockRewardInt)
		pool.Debt.Add(blockRewardInt)
		k.SetPool(ctx, pool)
	}

	pool := k.GetPool(ctx)
	if !pool.Debt.IsPositive() {
		return
	}

	fundAccAddr, err := sdk.AccAddressFromBech32(params.FundAddress)
	if err != nil {
		panic(err)
	}

	pay := pool.Debt
	balance := k.GetBalance(ctx, fundAccAddr, params.MintDenom)
	if balance.Amount.LT(pool.Debt) {
		pay = balance.Amount
	}

	if !pay.IsPositive() {
		return
	}

	payCoin := sdk.NewCoin(params.MintDenom, pay)
	err = k.AddCollectedFeesFromAccount(ctx, sdk.NewCoins(payCoin), fundAccAddr)
	if err != nil {
		panic(err)
	}

	pool.Debt = pool.Debt.Sub(pay)
	k.SetPool(ctx, pool)
}
