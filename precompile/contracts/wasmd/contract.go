package wasmd

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/precompile/contract"
	pcommon "github.com/kava-labs/kava/precompile/common"
)

type PrecompileExecutor struct {
	wasmdKeeper pcommon.WasmdKeeper
}

func (p PrecompileExecutor) instantiateCosmWasm(
	accessibleState contract.AccessibleState,
	caller common.Address,
	addr common.Address,
	packedInput []byte,
	suppliedGas uint64,
	readOnly bool,
	value *big.Int,
) (ret []byte, remainingGas uint64, err error) {
	return nil, 0, nil
}

func (p PrecompileExecutor) executeCosmWasm(
	accessibleState contract.AccessibleState,
	caller common.Address,
	addr common.Address,
	packedInput []byte,
	suppliedGas uint64,
	readOnly bool,
	value *big.Int,
) (ret []byte, remainingGas uint64, err error) {
	return nil, 0, nil
}

func (p PrecompileExecutor) queryCosmWasm(
	accessibleState contract.AccessibleState,
	caller common.Address,
	addr common.Address,
	packedInput []byte,
	suppliedGas uint64,
	readOnly bool,
	value *big.Int,
) (ret []byte, remainingGas uint64, err error) {
	return nil, 0, nil
}

// NewContract returns a new wasmd stateful precompiled contract.
//
//	This contract is used for testing purposes only and should not be used on public chains.
//	The functions of this contract (once implemented), will be used to exercise and test the various aspects of
//	the EVM such as gas usage, argument parsing, events, etc. The specific operations tested under this contract are
//	still to be determined.
func NewContract(wasmdKeeper pcommon.WasmdKeeper) (contract.StatefulPrecompiledContract, error) {

	executor := &PrecompileExecutor{
		wasmdKeeper: wasmdKeeper,
	}

	var functions []*contract.StatefulPrecompileFunction

	functions = append(functions, contract.NewStatefulPrecompileFunction(
		contract.MustCalculateFunctionSelector("instantiate(uint64,string,bytes,string,bytes)"),
		executor.instantiateCosmWasm,
	))

	functions = append(functions, contract.NewStatefulPrecompileFunction(
		contract.MustCalculateFunctionSelector("execute(string,bytes,bytes)"),
		executor.executeCosmWasm,
	))

	functions = append(functions, contract.NewStatefulPrecompileFunction(
		contract.MustCalculateFunctionSelector("query(string,bytes)"),
		executor.queryCosmWasm,
	))

	// Construct the contract with functions.
	precompile, err := contract.NewStatefulPrecompileContract(functions)

	if err != nil {
		return nil, fmt.Errorf("failed to instantiate wasmd precompile: %w", err)
	}

	return precompile, nil
}
