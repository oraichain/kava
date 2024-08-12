package cmd

import (
	"os"

	"cosmossdk.io/log"
	confixcmd "cosmossdk.io/tools/confix/cmd"
	tmcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	dbm "github.com/cosmos/cosmos-db"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	ethermintclient "github.com/evmos/ethermint/client"
	"github.com/evmos/ethermint/crypto/hd"
	servercfg "github.com/evmos/ethermint/server/config"
	ethermintflags "github.com/evmos/ethermint/server/flags"
	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/app/params"
	kavaclient "github.com/kava-labs/kava/client"
	"github.com/kava-labs/kava/migrate"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// EnvPrefix is the prefix environment variables must have to configure the app.
const EnvPrefix = "KAVA"

// NewRootCmd creates a new root command for the kava blockchain.
func NewRootCmd() *cobra.Command {
	app.SetSDKConfig().Seal()

	appOpts := simtestutil.NewAppOptionsWithFlagHome(tempDir())
	encodingConfig := app.MakeEncodingConfig()
	tempApp := app.NewApp(log.NewNopLogger(), dbm.NewMemDB(), cast.ToString(appOpts.Get(flags.FlagHome)), nil, encodingConfig, app.Options{
		SkipLoadLatest:        false,
		SkipGenesisInvariants: cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants)),
		InvariantCheckPeriod:  cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		MempoolEnableAuth:     false,
		EVMTrace:              cast.ToString(appOpts.Get(ethermintflags.EVMTracer)),
		EVMMaxGasWanted:       cast.ToUint64(appOpts.Get(ethermintflags.EVMMaxTxGasWanted)),
	})

	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithHomeDir(app.DefaultNodeHome).
		WithKeyringOptions(hd.EthSecp256k1Option()).
		WithViper(EnvPrefix)

	rootCmd := &cobra.Command{
		Use:   "kava",
		Short: "Daemon and CLI for the Kava blockchain.",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err = client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := servercfg.AppConfig("ukava")

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, tmcfg.DefaultConfig())
		},
	}

	addSubCmds(rootCmd, encodingConfig, app.DefaultNodeHome, *tempApp)

	return rootCmd
}

// addSubCmds registers all the sub commands used by kava.
func addSubCmds(rootCmd *cobra.Command, encodingConfig params.EncodingConfig, defaultNodeHome string, app app.App) {
	rootCmd.AddCommand(
		ethermintclient.ValidateChainID(
			genutilcli.InitCmd(app.BasicModuleManager, defaultNodeHome),
		),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, defaultNodeHome, genutiltypes.DefaultMessageValidator,
			encodingConfig.TxConfig.SigningContext().ValidatorAddressCodec()),
		migrate.MigrateGenesisCmd(),
		migrate.AssertInvariantsCmd(encodingConfig, app.BasicModuleManager),
		genutilcli.GenTxCmd(app.BasicModuleManager, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, defaultNodeHome, encodingConfig.TxConfig.SigningContext().ValidatorAddressCodec()),
		genutilcli.ValidateGenesisCmd(app.BasicModuleManager),
		AddGenesisAccountCmd(defaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true), // TODO add other shells, drop tmcli dependency, unhide?
		// testnetCmd(app.ModuleBasics, banktypes.GenesisBalancesIterator{}), // TODO add
		debug.Cmd(),
		confixcmd.ConfigCommand(),
	)

	ac := appCreator{
		encodingConfig: encodingConfig,
	}

	// ethermintserver adds additional flags to start the JSON-RPC server for evm support
	server.AddCommands(rootCmd, defaultNodeHome, ac.newApp, ac.appExport, ac.addStartCmdFlags)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		StatusCommand(),
		newQueryCmd(),
		newTxCmd(),
		kavaclient.KeyCommands(defaultNodeHome),
	)
}

var tempDir = func() string {
	dir, err := os.MkdirTemp("", "kavad")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}
	defer os.RemoveAll(dir)

	return dir
}
