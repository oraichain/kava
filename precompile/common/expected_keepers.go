package common

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type WasmdKeeper interface {
	Instantiate(ctx sdk.Context, codeID uint64, creator, admin sdk.AccAddress, initMsg []byte, label string, deposit sdk.Coins) (sdk.AccAddress, []byte, error)
	Execute(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, msg []byte, coins sdk.Coins) ([]byte, error)
	QuerySmart(ctx sdk.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error)
}
