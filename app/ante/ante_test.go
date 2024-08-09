package ante_test

import (
	"math/rand"
	"os"
	"testing"
	"time"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	tmdb "github.com/cosmos/cosmos-db"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	helpers "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"

	"github.com/kava-labs/kava/app"
)

func TestMain(m *testing.M) {
	app.SetSDKConfig()
	os.Exit(m.Run())
}

func TestAppAnteHandler_AuthorizedMempool(t *testing.T) {
	testPrivKeys, testAddresses := app.GeneratePrivKeyAddressPairs(10)
	unauthed := testAddresses[0:2]
	unauthedKeys := testPrivKeys[0:2]
	deputy := testAddresses[2]
	deputyKey := testPrivKeys[2]
	oracles := testAddresses[3:6]
	oraclesKeys := testPrivKeys[3:6]
	manual := testAddresses[6:]
	manualKeys := testPrivKeys[6:]

	encodingConfig := app.MakeEncodingConfig()

	opts := app.DefaultOptions
	opts.MempoolEnableAuth = true
	opts.MempoolAuthAddresses = manual

	tApp := app.TestApp{
		App: *app.NewApp(
			log.NewNopLogger(),
			tmdb.NewMemDB(),
			app.DefaultNodeHome,
			nil,
			encodingConfig,
			opts,
		),
	}

	chainID := "kavatest_1-1"
	tApp = tApp.InitializeFromGenesisStatesWithTimeAndChainID(
		time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
		chainID,
		app.NewFundedGenStateWithSameCoins(
			tApp.AppCodec(),
			sdk.NewCoins(sdk.NewInt64Coin("ukava", 1e9)),
			testAddresses,
		),
	)

	testcases := []struct {
		name       string
		address    sdk.AccAddress
		privKey    cryptotypes.PrivKey
		expectPass bool
	}{
		{
			name:       "unauthorized",
			address:    unauthed[1],
			privKey:    unauthedKeys[1],
			expectPass: false,
		},
		{
			name:       "oracle",
			address:    oracles[1],
			privKey:    oraclesKeys[1],
			expectPass: true,
		},
		{
			name:       "deputy",
			address:    deputy,
			privKey:    deputyKey,
			expectPass: true,
		},
		{
			name:       "manual",
			address:    manual[1],
			privKey:    manualKeys[1],
			expectPass: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			stdTx, err := helpers.GenSignedMockTx(
				rand.New(rand.NewSource(time.Now().UnixNano())),
				encodingConfig.TxConfig,
				[]sdk.Msg{
					banktypes.NewMsgSend(
						tc.address,
						testAddresses[0],
						sdk.NewCoins(sdk.NewInt64Coin("ukava", 1_000_000)),
					),
				},
				sdk.NewCoins(), // no fee
				helpers.DefaultGenTxGas,
				chainID,
				[]uint64{0},
				[]uint64{0}, // fixed sequence numbers will cause tests to fail sig verification if the same address is used twice
				tc.privKey,
			)
			require.NoError(t, err)
			txBytes, err := encodingConfig.TxConfig.TxEncoder()(stdTx)
			require.NoError(t, err)

			res, err := tApp.CheckTx(
				&abci.RequestCheckTx{
					Tx:   txBytes,
					Type: abci.CheckTxType_New,
				},
			)
			require.NoError(t, err)

			if tc.expectPass {
				require.Zero(t, res.Code, res.Log)
			} else {
				require.NotZero(t, res.Code)
			}
		})
	}
}

func TestAppAnteHandler_RejectMsgsInAuthz(t *testing.T) {
	testPrivKeys, testAddresses := app.GeneratePrivKeyAddressPairs(10)

	expiration := time.Date(9000, 1, 1, 0, 0, 0, 0, time.UTC)
	newMsgGrant := func(msgTypeUrl string) *authz.MsgGrant {
		msg, err := authz.NewMsgGrant(
			testAddresses[0],
			testAddresses[1],
			authz.NewGenericAuthorization(msgTypeUrl),
			&expiration,
		)
		if err != nil {
			panic(err)
		}
		return msg
	}

	chainID := "kavatest_1-1"
	encodingConfig := app.MakeEncodingConfig()

	testcases := []struct {
		name         string
		msg          sdk.Msg
		expectedCode uint32
	}{
		{
			name:         "MsgEthereumTx is blocked",
			msg:          newMsgGrant(sdk.MsgTypeURL(&evmtypes.MsgEthereumTx{})),
			expectedCode: sdkerrors.ErrUnauthorized.ABCICode(),
		},
		{
			name:         "MsgCreateVestingAccount is blocked",
			msg:          newMsgGrant(sdk.MsgTypeURL(&vestingtypes.MsgCreateVestingAccount{})),
			expectedCode: sdkerrors.ErrUnauthorized.ABCICode(),
		},
	}

	txEncoder := encodingConfig.TxConfig.TxEncoder()
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tApp := app.NewTestApp()

			tApp = tApp.InitializeFromGenesisStatesWithTimeAndChainID(
				time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
				chainID,
			)

			stdTx, err := helpers.GenSignedMockTx(
				rand.New(rand.NewSource(time.Now().UnixNano())),
				encodingConfig.TxConfig,
				[]sdk.Msg{tc.msg},
				sdk.NewCoins(), // no fee
				helpers.DefaultGenTxGas,
				chainID,
				[]uint64{0},
				[]uint64{0},
				testPrivKeys[0],
			)
			require.NoError(t, err)
			txBytes, err := txEncoder(stdTx)
			require.NoError(t, err)

			resCheckTx, err := tApp.CheckTx(
				&abci.RequestCheckTx{
					Tx:   txBytes,
					Type: abci.CheckTxType_New,
				},
			)
			require.NoError(t, err)
			require.Equal(t, resCheckTx.Code, tc.expectedCode, resCheckTx.Log)

			_, result, err := tApp.SimDeliver(
				txEncoder,
				stdTx,
			)
			_, code, _ := errorsmod.ABCIInfo(err, false)
			require.Equal(t, code, tc.expectedCode, result.Log)
		})
	}
}
