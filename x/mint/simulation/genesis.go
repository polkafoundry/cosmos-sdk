package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
)

// Simulation parameter constants
const (
	Inflation           = "inflation"
	InflationRateChange = "inflation_rate_change"
	InflationMax        = "inflation_max"
	InflationMin        = "inflation_min"
	GoalBonded          = "goal_bonded"
	ValidatorRewardFund = "validator_reward_fund"
	FundAddress         = "fund_address"
)

// GenInflation randomized Inflation
func GenInflation(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(r.Intn(99)), 2)
}

// GenInflationRateChange randomized InflationRateChange
func GenInflationRateChange(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(r.Intn(99)), 2)
}

// GenInflationMax randomized InflationMax
func GenInflationMax(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(20, 2)
}

// GenInflationMin randomized InflationMin
func GenInflationMin(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(7, 2)
}

// GenGoalBonded randomized GoalBonded
func GenGoalBonded(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(67, 2)
}

func GenValidatorRewardFund(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(r.Uint64()), 2)
}

func GenFundAddress(r *rand.Rand) string {
	seed := []byte(strconv.Itoa(r.Int()))
	pubKey := ed25519.GenPrivKeyFromSecret(seed).PubKey()
	return sdk.AccAddress(pubKey.Address()).String()
}

// RandomizedGenState generates a random GenesisState for mint
func RandomizedGenState(simState *module.SimulationState) {
	// minter
	var inflation sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, Inflation, &inflation, simState.Rand,
		func(r *rand.Rand) { inflation = GenInflation(r) },
	)

	// params
	var inflationRateChange sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, InflationRateChange, &inflationRateChange, simState.Rand,
		func(r *rand.Rand) { inflationRateChange = GenInflationRateChange(r) },
	)

	var inflationMax sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, InflationMax, &inflationMax, simState.Rand,
		func(r *rand.Rand) { inflationMax = GenInflationMax(r) },
	)

	var inflationMin sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, InflationMin, &inflationMin, simState.Rand,
		func(r *rand.Rand) { inflationMin = GenInflationMin(r) },
	)

	var goalBonded sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, GoalBonded, &goalBonded, simState.Rand,
		func(r *rand.Rand) { goalBonded = GenGoalBonded(r) },
	)

	mintDenom := sdk.DefaultBondDenom
	blocksPerYear := uint64(60 * 60 * 8766 / 5)

	var validatorRewardFund sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, ValidatorRewardFund, &validatorRewardFund, simState.Rand,
		func(r *rand.Rand) { validatorRewardFund = GenValidatorRewardFund(r) },
	)

	fundDuration := 5 * blocksPerYear

	var fundAddress string
	simState.AppParams.GetOrGenerate(
		simState.Cdc, FundAddress, &fundAddress, simState.Rand,
		func(r *rand.Rand) { fundAddress = GenFundAddress(r) },
	)

	params := types.NewParams(mintDenom, inflationRateChange, inflationMax, inflationMin, goalBonded, blocksPerYear,
		false, validatorRewardFund, fundDuration, fundAddress,
	)

	mintGenesis := types.NewGenesisState(types.InitialMinter(inflation), params, types.DefaultInitialPool())

	bz, err := json.MarshalIndent(&mintGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated minting parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(mintGenesis)
}
