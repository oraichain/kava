package app

import (
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/gorilla/mux"
)

// RegisterLegacyTxRoutes registers a legacy tx routes that use amino encoding json
func RegisterLegacyTxRoutes(clientCtx client.Context, r *mux.Router) {

}

// LegacyTxBroadcastRequest represents a broadcast request with an amino json encoded transaction
type LegacyTxBroadcastRequest struct {
	Tx   legacytx.StdTx `json:"tx"`
	Mode string         `json:"mode"`
}

var _ codectypes.UnpackInterfacesMessage = LegacyTxBroadcastRequest{}

// UnpackInterfaces implements the UnpackInterfacesMessage interface
func (m LegacyTxBroadcastRequest) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return m.Tx.UnpackInterfaces(unpacker)
}
