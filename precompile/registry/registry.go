package registry

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/precompile/contract"
	"github.com/ethereum/go-ethereum/precompile/modules"
	pcommon "github.com/kava-labs/kava/precompile/common"
	"github.com/kava-labs/kava/precompile/contracts/wasmd"
)

const (
	// WasmdContractAddress the primary noop contract address for testing
	WasmdContractAddress = "0x9000000000000000000000000000000000000001"
)

// init registers stateful precompile contracts with the global precompile registry
// defined in kava-labs/go-ethereum/precompile/modules
func InitializePrecompiles(wasmdKeeper pcommon.WasmdKeeper, wasmdViewKeeper pcommon.WasmdViewKeeper, evmKeeper pcommon.EVMKeeper) {
	wasmdContract, err := wasmd.NewContract(wasmdKeeper, wasmdViewKeeper, evmKeeper)
	if err != nil {
		panic(fmt.Errorf("error creating contract for address %s: %w", WasmdContractAddress, err))
	}
	register(WasmdContractAddress, wasmdContract)
}

// register accepts a 0x address string and a stateful precompile contract constructor, instantiates the
// precompile contract via the constructor, and registers it with the precompile module registry.
//
// This panics if the contract can not be created or the module can not be registered
func register(address string, contract contract.StatefulPrecompiledContract) {
	module := modules.Module{
		Address:  common.HexToAddress(address),
		Contract: contract,
	}

	err := modules.RegisterModule(module)
	if err != nil {
		panic(fmt.Errorf("error registering contract module for address %s: %w", address, err))
	}
}
