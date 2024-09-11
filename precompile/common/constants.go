package common

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// OraiToSoraiMultiplier Fields that were denominated in orai will be converted to sorai (1orai = 10^12sorai)
// for existing Ethereum application (which assumes 18 decimal points) to display properly.
var oraiToSoraiMultiplier = big.NewInt(1_000_000_000_000)
var SdkOraiToSoraiMultiplier = sdk.NewIntFromBigInt(oraiToSoraiMultiplier)
