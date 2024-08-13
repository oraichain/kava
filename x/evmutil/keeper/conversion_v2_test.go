package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	"github.com/kava-labs/kava/x/evmutil/testutil"
	"github.com/kava-labs/kava/x/evmutil/types"
)

type ConversionV2TestSuite struct {
	testutil.NewSuite
}

func TestConversionV2TestSuite(t *testing.T) {
	suite.Run(t, new(ConversionV2TestSuite))
}

func (suite *ConversionV2TestSuite) TestConvertCoinToERC20() {
	contractAddr := suite.DeployERC20()

	pair := types.NewConversionPair(
		contractAddr,
		"erc20/usdc",
	)

	amount := big.NewInt(100)
	originAcc := sdk.AccAddress(suite.Key1.PubKey().Address().Bytes())
	recipientAcc := types.NewInternalEVMAddress(common.BytesToAddress(suite.Key2.PubKey().Address()))
	moduleAddr := types.NewInternalEVMAddress(types.ModuleEVMAddress)

	// Starting balance of origin account
	coin, err := suite.Keeper.MintConversionPairCoin(suite.Ctx, pair, amount, originAcc)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.NewCoin(pair.Denom, sdkmath.NewIntFromBigInt(amount)), coin)

	// Mint same initial balance for module account as backing erc20 supply
	err = suite.Keeper.MintERC20(
		suite.Ctx,
		pair.GetAddress(), // contractAddr
		moduleAddr,        //receiver
		amount,
	)
	suite.Require().NoError(err)

	startBal, err := suite.Keeper.QueryERC20Supply(suite.Ctx, contractAddr)
	fmt.Println("total supply: ", startBal, err)
	suite.Require().NoError(err)

	// convert coin to erc20
	ctx := suite.Ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
	err = suite.Keeper.ConvertCoinToERC20(
		ctx,
		originAcc,
		recipientAcc,
		sdk.NewCoin(pair.Denom, sdkmath.NewIntFromBigInt(amount)),
	)
	suite.Require().NoError(err)
	suite.Require().LessOrEqual(ctx.GasMeter().GasConsumed(), uint64(500000))
	suite.Require().GreaterOrEqual(ctx.GasMeter().GasConsumed(), uint64(50000))

	// Source should decrease
	bal := suite.App.GetBankKeeper().GetBalance(suite.Ctx, originAcc, pair.Denom)
	suite.Require().Equal(sdkmath.ZeroInt(), bal.Amount, "conversion should decrease source balance")

	// Module bal should also decrease
	moduleBal := suite.GetERC20BalanceOf(
		types.ERC20MintableBurnableContract.ABI,
		pair.GetAddress(),
		moduleAddr,
	)
	suite.Require().Equal(
		// String() due to non-equal struct values for 0
		big.NewInt(0).String(),
		moduleBal.String(),
		"balance should decrease module account by unlock amount",
	)

	// Recipient balance should increase by same amount
	recipientBal := suite.GetERC20BalanceOf(
		types.ERC20MintableBurnableContract.ABI,
		pair.GetAddress(),
		recipientAcc,
	)
	suite.Require().Equal(
		// String() due to non-equal struct values for 0
		amount,
		recipientBal,
		"recipient balance should increase",
	)

	suite.EventsContains(suite.GetEvents(),
		sdk.NewEvent(
			types.EventTypeConvertCoinToERC20,
			sdk.NewAttribute(types.AttributeKeyInitiator, originAcc.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, recipientAcc.String()),
			sdk.NewAttribute(types.AttributeKeyERC20Address, pair.GetAddress().String()),
			sdk.NewAttribute(types.AttributeKeyAmount, coin.String()),
		))
}
