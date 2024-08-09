package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

var _ govtypes.TallyHandler = TallyHandler{}

// TallyHandler is the tally handler for kava
type TallyHandler struct {
	gk  govkeeper.Keeper
	stk stakingkeeper.Keeper

	bk bankkeeper.Keeper
}

// NewTallyHandler creates a new tally handler.
func NewTallyHandler(
	gk govkeeper.Keeper, stk stakingkeeper.Keeper,
	bk bankkeeper.Keeper,
) TallyHandler {
	return TallyHandler{
		gk:  gk,
		stk: stk,

		bk: bk,
	}
}

func (th TallyHandler) Tally(ctx sdk.Context, proposal types.Proposal) (passes bool, burnDeposits bool, tallyResults types.TallyResult) {
	results := make(map[types.VoteOption]sdk.Dec)
	results[types.OptionYes] = sdk.ZeroDec()
	results[types.OptionAbstain] = sdk.ZeroDec()
	results[types.OptionNo] = sdk.ZeroDec()
	results[types.OptionNoWithVeto] = sdk.ZeroDec()

	totalVotingPower := sdk.ZeroDec()
	currValidators := make(map[string]types.ValidatorGovInfo)

	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		if len(val.Vote) == 0 {
			continue
		}

		sharesAfterDeductions := val.DelegatorShares.Sub(val.DelegatorDeductions)
		votingPower := sharesAfterDeductions.MulInt(val.BondedTokens).Quo(val.DelegatorShares)

		for _, option := range val.Vote {
			subPower := votingPower.Mul(option.Weight)
			results[option.Option] = results[option.Option].Add(subPower)
		}
		totalVotingPower = totalVotingPower.Add(votingPower)
	}

	tallyParams := th.gk.GetTallyParams(ctx)
	tallyResults = types.NewTallyResultFromMap(results)

	// TODO: Upgrade the spec to cover all of these cases & remove pseudocode.
	// If there is no staked coins, the proposal fails
	if th.stk.TotalBondedTokens(ctx).IsZero() {
		return false, false, tallyResults
	}

	// If there is not enough quorum of votes, the proposal fails
	percentVoting := totalVotingPower.Quo(th.stk.TotalBondedTokens(ctx).ToDec())
	if percentVoting.LT(tallyParams.Quorum) {
		return false, true, tallyResults
	}

	// If no one votes (everyone abstains), proposal fails
	if totalVotingPower.Sub(results[types.OptionAbstain]).Equal(sdk.ZeroDec()) {
		return false, false, tallyResults
	}

	// If more than 1/3 of voters veto, proposal fails
	if results[types.OptionNoWithVeto].Quo(totalVotingPower).GT(tallyParams.VetoThreshold) {
		return false, true, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes
	if results[types.OptionYes].Quo(totalVotingPower.Sub(results[types.OptionAbstain])).GT(tallyParams.Threshold) {
		return true, false, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote No, proposal fails
	return false, false, tallyResults
}

// bkavaByDenom a map of the bkava denom and the amount of bkava for that denom.
type bkavaByDenom map[string]sdk.Int

func (bkavaMap bkavaByDenom) add(coin sdk.Coin) {
	_, found := bkavaMap[coin.Denom]
	if !found {
		bkavaMap[coin.Denom] = sdk.ZeroInt()
	}
	bkavaMap[coin.Denom] = bkavaMap[coin.Denom].Add(coin.Amount)
}

func (bkavaMap bkavaByDenom) toCoins() sdk.Coins {
	coins := sdk.Coins{}
	for denom, amt := range bkavaMap {
		coins = coins.Add(sdk.NewCoin(denom, amt))
	}
	return coins.Sort()
}
