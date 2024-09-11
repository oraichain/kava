package wasmd_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/precompile/contracts/wasmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockWasmer struct {
}

func (m *MockWasmer) Instantiate(ctx sdk.Context, codeID uint64, creator, admin sdk.AccAddress, initMsg []byte, label string, deposit sdk.Coins) (sdk.AccAddress, []byte, error) {
	addr := sdk.MustAccAddressFromBech32("orai19xtunzaq20unp8squpmfrw8duclac22hd7ves2")
	return addr, []byte("Instantiate"), nil
}

func (m *MockWasmer) Execute(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, msg []byte, coins sdk.Coins) ([]byte, error) {
	return []byte("Execute"), nil
}

func (m *MockWasmer) QuerySmart(ctx sdk.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error) {
	return []byte("QuerySmart"), nil
}

// TestContractConstructor ensures we have a valid constructor. This will fail
// if we attempt to define invalid or duplicate function selectors.
func TestContractConstructor(t *testing.T) {
	wasmer := &MockWasmer{}
	precompile, err := wasmd.NewContract(wasmer)
	require.NoError(t, err, "expected precompile not error when created")
	assert.NotNil(t, precompile, "expected precompile contract to be defined")
}
