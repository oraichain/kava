package app

import (
	"context"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

const (
	MainnetUpgradeName = "v0.21.0"
	TestnetUpgradeName = "v0.21.0-alpha.0"
)

func (app App) RegisterUpgradeHandlers(db dbm.DB) {
	// register upgrade handler for mainnet
	app.upgradeKeeper.SetUpgradeHandler(MainnetUpgradeName, MainnetUpgradeHandler(app))
	// register upgrade handler for testnet
	app.upgradeKeeper.SetUpgradeHandler(TestnetUpgradeName, TestnetUpgradeHandler(app))

	upgradeInfo, err := app.upgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	// MAINNET STORE CHANGES
	// only the community module is added which has no store.
	// therefore, no store migration is necessary for mainnet.

	// TESTNET STORE CHANGES
	// we must undo the store changes performed in the v0.20.0-alpha.0 upgrade handler.
	if upgradeInfo.Name == TestnetUpgradeName && !app.upgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{
				minttypes.StoreKey,
			},
			Deleted: []string{
				"kavamint",
			},
		}
		// override the store loader to handle cleaning up bad testnet x/mint state
		app.SetStoreLoader(TestnetStoreLoader(app, db, upgradeInfo.Height, &storeUpgrades))
	}
}

// TestnetStoreLoader removes the previous iavl tree for the mint module, ensuring even store heights without
// modifications to iavl to support non-consecutive versions and deletion of all nodes for a new tree at the upgrade height
func TestnetStoreLoader(app App, db dbm.DB, upgradeHeight int64, storeUpgrades *storetypes.StoreUpgrades) baseapp.StoreLoader {
	return func(ms storetypes.CommitMultiStore) error {
		// if this is the upgrade height, delete all remnant x/mint store versions to ensure we start from clean slate
		if upgradeHeight == ms.LastCommitID().Version+1 {
			app.Logger().Info("removing x/mint historic versions from store")
			prefix := "s/k:" + minttypes.StoreKey + "/"

			// The mint module iavl versioned tree is stored at "s/k:mint/"
			prefixdb := dbm.NewPrefixDB(db, []byte(prefix))

			itr, err := prefixdb.Iterator(nil, nil)
			if err != nil {
				return err
			}

			// Collect keys since deletion during iteration may cause issues
			var keys [][]byte
			for itr.Valid() {
				keys = append(keys, itr.Key())
				itr.Next()
			}
			itr.Close()

			// Delete all keys and thus all history of the mint store iavl tree
			for _, k := range keys {
				prefixdb.Delete(k)
			}
		}

		// run the standard upgrade handler, now starting at a clean state for the mint store key
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		return upgradetypes.UpgradeStoreLoader(upgradeHeight, storeUpgrades)(ms)
	}
}

// MainnetUpgradeHandler does nothing. No state changes are necessary on mainnet because v0.20.0 was
// never released and its upgrade handler was never run.
func MainnetUpgradeHandler(app App) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// no-op
		app.Logger().Info("running mainnet upgrade handler")
		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	}
}

// TestnetUpgradeHandler is the inverse of the v0.20.0-alpha.0 upgrade handler that was run on public
// testnet. It reverts the state changes to bring the chain back to its v0.19.0 state, which is expected
// in this upgrade.
// See original handler here: https://github.com/Kava-Labs/kava/blob/v0.20.0-alpha.0/app/upgrades.go#L65
func TestnetUpgradeHandler(app App) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		app.Logger().Info("running testnet upgrade handler")

		// move community pool funds back to community pool from community module.
		app.Logger().Info("migrating community pool funds")

		// reenable community tax
		app.Logger().Info("re-enabling community tax")
		ReenableCommunityTax(ctx, app.distrKeeper)

		// remove mint from the version map to ensure InitGenesis for x/mint is run
		delete(fromVM, "mint")

		// run migrations
		vm, err := app.mm.RunMigrations(ctx, app.configurator, fromVM)
		if err != nil {
			panic(err)
		}

		// initialize x/mint params. must be done after migrations so module exists.
		app.Logger().Info("initializing x/mint state")
		InitializeMintState(ctx, app.mintKeeper, app.stakingKeeper)

		return vm, nil
	}
}

// ReenableCommunityTax sets x/distribution's community_tax to the value currently on mainnet.
func ReenableCommunityTax(ctx context.Context, distrKeeper distrkeeper.Keeper) {
	params, err := distrKeeper.Params.Get(ctx)
	if err != nil {
		panic(err)
	}
	params.CommunityTax = sdkmath.LegacyMustNewDecFromStr("0.925000000000000000") // community tax currently present on mainnet
	distrKeeper.Params.Set(ctx, params)
}

// InitializeMintState sets up the parameters and state of x/mint.
func InitializeMintState(
	ctx context.Context,
	mintKeeper mintkeeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
) {
	var err error
	// init params for x/mint with values from mainnet
	inflationRate := sdkmath.LegacyMustNewDecFromStr("0.750000000000000000")
	params := minttypes.DefaultParams()
	params.MintDenom, err = stakingKeeper.BondDenom(ctx)
	if err != nil {
		panic(err)
	}
	params.InflationMax = inflationRate
	params.InflationMin = inflationRate

	mintKeeper.Params.Set(ctx, params)
}
