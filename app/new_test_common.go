package app

import (
	"encoding/json"
	"testing"
	"time"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmdb "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	distkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"

	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	evmutilkeeper "github.com/kava-labs/kava/x/evmutil/keeper"

	"github.com/stretchr/testify/require"
)

// TestAppV2 is a simple wrapper around an App. It exposes internal keepers for use in integration tests.
// This file also contains test helpers. Ideally they would be in separate package.
// Basic Usage:
//
//	Create a test app with NewTestApp, then all keepers and their methods can be accessed for test setup and execution.
//
// Advanced Usage:
//
//	Some tests call for an app to be initialized with some state. This can be achieved through keeper method calls (ie keeper.SetParams(...)).
//	However this leads to a lot of duplicated logic similar to InitGenesis methods.
//	So TestAppV2.InitializeFromGenesisStates() will call InitGenesis with the default genesis state.
//	and TestAppV2.InitializeFromGenesisStates(authState, cdpState) will do the same but overwrite the auth and cdp sections of the default genesis state
//	Creating the genesis states can be combersome, but helper methods can make it easier such as NewAuthGenStateFromAccounts below.
type TestAppV2 struct {
	App
}

// NewTestApp creates a new TestAppV2
//
// Note, it also sets the sdk config with the app's address prefix, coin type, etc.
func NewTestAppV2() TestAppV2 {
	SetSDKConfig()

	return NewTestAppV2FromSealed()
}

// NewTestAppFromSealed creates a TestAppV2 without first setting sdk config.
func NewTestAppV2FromSealed() TestAppV2 {
	db := tmdb.NewMemDB()

	encCfg := MakeEncodingConfig()

	app := NewApp(log.NewNopLogger(), db, DefaultNodeHome, nil, encCfg, DefaultOptions, baseapp.SetChainID(TestChainID))
	cdc := app.AppCodec()
	genesisState := app.NewTestGenesisState(cdc)
	// genesisState[types.ModuleName] = cdc.MustMarshalJSON(types.DefaultGenesisState())
	evmGenesis := evmtypes.DefaultGenesisState()
	evmGenesis.Params.EvmDenom = "akava"

	feemarketGenesis := feemarkettypes.DefaultGenesisState()
	feemarketGenesis.Params.EnableHeight = 1
	feemarketGenesis.Params.NoBaseFee = false
	genesisState[evmtypes.ModuleName] = cdc.MustMarshalJSON(evmGenesis)
	genesisState[feemarkettypes.ModuleName] = cdc.MustMarshalJSON(feemarketGenesis)

	stateBytes, _ := json.MarshalIndent(genesisState, "", "  ")

	_, err := app.InitChain(
		&abci.RequestInitChain{
			Time:          emptyTime,
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
			ChainId:       TestChainID,
			// Set consensus params, which is needed by x/feemarket
			ConsensusParams: &cmtproto.ConsensusParams{
				Block: &cmtproto.BlockParams{
					MaxBytes: 200000,
					MaxGas:   20000000,
				},
			},
			InitialHeight: defaultInitialHeight,
		},
	)
	if err != nil {
		panic(err)
	}
	// _, err = app.Commit()
	// if err != nil {
	// 	panic(err)
	// }
	// _, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
	// 	Height: app.LastBlockHeight() + 1, Time: emptyTime,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	return TestAppV2{App: *app}
}

// nolint
func (tApp TestAppV2) GetAccountKeeper() authkeeper.AccountKeeper { return tApp.accountKeeper }
func (tApp TestAppV2) GetBankKeeper() bankkeeper.Keeper           { return tApp.bankKeeper }
func (tApp TestAppV2) GetMintKeeper() mintkeeper.Keeper           { return tApp.mintKeeper }
func (tApp TestAppV2) GetStakingKeeper() *stakingkeeper.Keeper    { return tApp.stakingKeeper }
func (tApp TestAppV2) GetSlashingKeeper() slashingkeeper.Keeper   { return tApp.slashingKeeper }
func (tApp TestAppV2) GetDistrKeeper() distkeeper.Keeper          { return tApp.distrKeeper }
func (tApp TestAppV2) GetGovKeeper() *govkeeper.Keeper            { return tApp.govKeeper }
func (tApp TestAppV2) GetCrisisKeeper() *crisiskeeper.Keeper      { return tApp.crisisKeeper }
func (tApp TestAppV2) GetParamsKeeper() paramskeeper.Keeper       { return tApp.paramsKeeper }

func (tApp TestAppV2) GetEvmutilKeeper() evmutilkeeper.Keeper { return tApp.evmutilKeeper }
func (tApp TestAppV2) GetEvmKeeper() *evmkeeper.Keeper        { return tApp.evmKeeper }

func (tApp TestAppV2) GetFeeMarketKeeper() feemarketkeeper.Keeper { return tApp.feeMarketKeeper }

// InitializeFromGenesisStates calls InitChain on the app using the provided genesis states.
// If any module genesis states are missing, defaults are used.
func (tApp TestAppV2) InitializeFromGenesisStates(genesisStates ...GenesisState) TestAppV2 {
	return tApp.InitializeFromGenesisStatesWithTimeAndChainIDAndHeight(emptyTime, TestChainID, defaultInitialHeight, genesisStates...)
}

// InitializeFromGenesisStatesWithTime calls InitChain on the app using the provided genesis states and time.
// If any module genesis states are missing, defaults are used.
func (tApp TestAppV2) InitializeFromGenesisStatesWithTime(genTime time.Time, genesisStates ...GenesisState) TestAppV2 {
	return tApp.InitializeFromGenesisStatesWithTimeAndChainIDAndHeight(genTime, TestChainID, defaultInitialHeight, genesisStates...)
}

// InitializeFromGenesisStatesWithTimeAndChainID calls InitChain on the app using the provided genesis states, time, and chain id.
// If any module genesis states are missing, defaults are used.
func (tApp TestAppV2) InitializeFromGenesisStatesWithTimeAndChainID(genTime time.Time, chainID string, genesisStates ...GenesisState) TestAppV2 {
	return tApp.InitializeFromGenesisStatesWithTimeAndChainIDAndHeight(genTime, chainID, defaultInitialHeight, genesisStates...)
}

// InitializeFromGenesisStatesWithTimeAndChainIDAndHeight calls InitChain on the app using the provided genesis states and other parameters.
// If any module genesis states are missing, defaults are used.
func (tApp TestAppV2) InitializeFromGenesisStatesWithTimeAndChainIDAndHeight(genTime time.Time, chainID string, initialHeight int64, genesisStates ...GenesisState) TestAppV2 {
	// Create a default genesis state and overwrite with provided values
	genesisState := tApp.NewDefaultGenesisState()
	for _, state := range genesisStates {
		for k, v := range state {
			genesisState[k] = v
		}
	}

	// Initialize the chain
	stateBytes, err := json.Marshal(genesisState)
	if err != nil {
		panic(err)
	}
	_, err = tApp.InitChain(
		&abci.RequestInitChain{
			Time:          genTime,
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
			ChainId:       chainID,
			// Set consensus params, which is needed by x/feemarket
			ConsensusParams: &cmtproto.ConsensusParams{
				Block: &cmtproto.BlockParams{
					MaxBytes: 200000,
					MaxGas:   20000000,
				},
			},
			InitialHeight: initialHeight,
		},
	)
	tApp.Commit()
	tApp.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height: tApp.LastBlockHeight() + 1, Time: genTime,
	})
	return tApp
}

// CheckBalance requires the account address has the expected amount of coins.
func (tApp TestAppV2) CheckBalance(t *testing.T, ctx sdk.Context, owner sdk.AccAddress, expectedCoins sdk.Coins) {
	coins := tApp.GetBankKeeper().GetAllBalances(ctx, owner)
	require.Equal(t, expectedCoins, coins)
}

// GetModuleAccountBalance gets the current balance of the denom for a module account
func (tApp TestAppV2) GetModuleAccountBalance(ctx sdk.Context, moduleName string, denom string) sdkmath.Int {
	moduleAcc := tApp.accountKeeper.GetModuleAccount(ctx, moduleName)
	balance := tApp.bankKeeper.GetBalance(ctx, moduleAcc.GetAddress(), denom)
	return balance.Amount
}

// FundAccount is a utility function that funds an account by minting and sending the coins to the address.
func (tApp TestAppV2) FundAccount(ctx sdk.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := tApp.bankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}

	return tApp.bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, amounts)
}

// NewQueryServerTestHelper creates a new QueryServiceTestHelper that wraps the provided sdk.Context.
func (tApp TestAppV2) NewQueryServerTestHelper(ctx sdk.Context) *baseapp.QueryServiceTestHelper {
	return baseapp.NewQueryServerTestHelper(ctx, tApp.interfaceRegistry)
}

// FundModuleAccount is a utility function that funds a module account by minting and sending the coins to the address.
func (tApp TestAppV2) FundModuleAccount(ctx sdk.Context, recipientMod string, amounts sdk.Coins) error {
	if err := tApp.bankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}

	return tApp.bankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, recipientMod, amounts)
}

// CreateNewUnbondedValidator creates a new validator in the staking module.
// New validators are unbonded until the end blocker is run.
func (tApp TestAppV2) CreateNewUnbondedValidator(ctx sdk.Context, valAddress sdk.ValAddress, selfDelegation sdkmath.Int) error {
	denom, err := tApp.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return err
	}
	msg, err := stakingtypes.NewMsgCreateValidator(
		valAddress.String(),
		ed25519.GenPrivKey().PubKey(),
		sdk.NewCoin(denom, selfDelegation),
		stakingtypes.Description{},
		stakingtypes.NewCommissionRates(sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec()),
		sdkmath.NewInt(1e6),
	)
	if err != nil {
		return err
	}

	msgServer := stakingkeeper.NewMsgServerImpl(tApp.stakingKeeper)
	_, err = msgServer.CreateValidator(sdk.WrapSDKContext(ctx), msg)
	return err
}
