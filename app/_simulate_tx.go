package app

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/auth"
)

// SimulateRequest represents attributes of a tx that will be simulated
type SimulateRequest struct {
	Msgs []sdk.Msg   `json:"msgs"`
	Fee  auth.StdFee `json:"fee"`
	Memo string      `json:"memo"`
}

// RegisterSimulateRoutes registers a tx simulate route to a mux router with
// a provided cli context
func RegisterSimulateRoutes(cliCtx context.CLIContext, r *mux.Router) {

}
