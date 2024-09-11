package wasmd

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/precompile/contract"
	pcommon "github.com/kava-labs/kava/precompile/common"
)

// Singleton StatefulPrecompiledContract.
var (
	// IBCRawABI contains the raw ABI of IBC contract.
	//go:embed abi.json
	IBCRawABI string

	IBCABI = contract.MustParseABI(IBCRawABI)
)

type PrecompileExecutor struct {
	wasmdKeeper pcommon.WasmdKeeper
	evmKeeper   pcommon.EVMKeeper
}

func (p PrecompileExecutor) instantiateCosmWasm(
	accessibleState contract.AccessibleState,
	caller common.Address,
	callingContract common.Address,
	packedInput []byte,
	suppliedGas uint64,
	readOnly bool,
	value *big.Int,
) (ret []byte, remainingGas uint64, rerr error) {

	defer func() {
		if err := recover(); err != nil {
			ret = nil
			remainingGas = 0
			rerr = fmt.Errorf("%s", err)
			return
		}
	}()
	if readOnly {
		rerr = errors.New("cannot call instantiate from staticcall")
		return
	}

	if !bytes.Equal(caller.Bytes(), callingContract.Bytes()) {
		rerr = errors.New("cannot delegatecall instantiate")
		return
	}

	res, err := IBCABI.UnpackInput("instantiate", packedInput)
	if err != nil {
		rerr = err
		return
	}

	codeID := *abi.ConvertType(res[0], new(uint64)).(*uint64)
	admin := *abi.ConvertType(res[1], new(string)).(*string)
	msg := *abi.ConvertType(res[2], new([]byte)).(*[]byte)
	label := *abi.ConvertType(res[3], new(string)).(*string)

	ctxer, ok := accessibleState.GetStateDB().(pcommon.Contexter)
	if !ok {
		rerr = errors.New("cannot get context from EVM")
		return
	}
	ctx := ctxer.Ctx()

	creator := p.evmKeeper.GetCosmosAddressMapping(ctx, caller)

	baseDenom, err := sdk.GetBaseDenom()
	if err != nil {
		rerr = err
		return
	}
	coinsValue := sdk.NewIntFromBigInt(value).Quo(pcommon.SdkOraiToSoraiMultiplier)
	deposit := sdk.NewCoins(sdk.NewCoin(baseDenom, coinsValue))

	adminAddr, err := sdk.AccAddressFromBech32(admin)
	if err != nil {
		rerr = err
		return
	}

	addr, data, err := p.wasmdKeeper.Instantiate(ctx, codeID, creator, adminAddr, msg, label, deposit)
	if err != nil {
		rerr = err
		return
	}

	cosmosGasUsed := ctx.GasMeter().GasConsumed()

	ret, rerr = IBCABI.Pack("instantiate", addr.String(), data)

	remainingGas, rerr = contract.DeductGas(suppliedGas, cosmosGasUsed)

	return
}

func (p PrecompileExecutor) executeCosmWasm(
	accessibleState contract.AccessibleState,
	caller common.Address,
	callingContract common.Address,
	packedInput []byte,
	suppliedGas uint64,
	readOnly bool,
	value *big.Int,
) (ret []byte, remainingGas uint64, rerr error) {

	defer func() {
		if err := recover(); err != nil {
			ret = nil
			remainingGas = 0
			rerr = fmt.Errorf("%s", err)
			return
		}
	}()
	if readOnly {
		rerr = errors.New("cannot call execute from staticcall")
		return
	}

	res, err := IBCABI.UnpackInput("execute", packedInput)
	if err != nil {
		rerr = err
		return
	}

	contractAddress := *abi.ConvertType(res[0], new(string)).(*string)
	msg := *abi.ConvertType(res[1], new([]byte)).(*[]byte)

	ctxer, ok := accessibleState.GetStateDB().(pcommon.Contexter)
	if !ok {
		rerr = errors.New("cannot get context from EVM")
		return
	}
	ctx := ctxer.Ctx()

	senderAddr := p.evmKeeper.GetCosmosAddressMapping(ctx, caller)

	baseDenom, err := sdk.GetBaseDenom()
	if err != nil {
		rerr = err
		return
	}
	coinsValue := sdk.NewIntFromBigInt(value).Quo(pcommon.SdkOraiToSoraiMultiplier)
	deposit := sdk.NewCoins(sdk.NewCoin(baseDenom, coinsValue))

	// addresses will be sent in Cosmos format
	contractAddr, err := sdk.AccAddressFromBech32(contractAddress)
	if err != nil {
		rerr = err
		return
	}

	exeRes, err := p.wasmdKeeper.Execute(ctx, contractAddr, senderAddr, msg, deposit)

	if err != nil {
		rerr = err
		return
	}

	cosmosGasUsed := ctx.GasMeter().GasConsumed()

	ret, rerr = IBCABI.Pack("execute", exeRes)

	remainingGas, rerr = contract.DeductGas(suppliedGas, cosmosGasUsed)

	return

}

func (p PrecompileExecutor) queryCosmWasm(
	accessibleState contract.AccessibleState,
	caller common.Address,
	addr common.Address,
	packedInput []byte,
	suppliedGas uint64,
	readOnly bool,
	value *big.Int,
) (ret []byte, remainingGas uint64, rerr error) {

	defer func() {
		if err := recover(); err != nil {
			ret = nil
			remainingGas = 0
			rerr = fmt.Errorf("%s", err)
			return
		}
	}()
	if value != nil && value.Sign() != 0 {
		rerr = errors.New("sending funds to a non-payable function")
		return
	}

	res, err := IBCABI.UnpackInput("query", packedInput)
	if err != nil {
		rerr = err
		return
	}

	contractAddress := *abi.ConvertType(res[0], new(string)).(*string)
	req := *abi.ConvertType(res[1], new([]byte)).(*[]byte)

	ctxer, ok := accessibleState.GetStateDB().(pcommon.Contexter)
	if !ok {
		rerr = errors.New("cannot get context from EVM")
		return
	}
	ctx := ctxer.Ctx()

	// addresses will be sent in Cosmos format
	contractAddr, err := sdk.AccAddressFromBech32(contractAddress)
	if err != nil {
		rerr = err
		return
	}

	queryRes, err := p.wasmdKeeper.QuerySmart(ctx, contractAddr, req)
	if err != nil {
		rerr = err
		return
	}

	cosmosGasUsed := ctx.GasMeter().GasConsumed()

	ret, rerr = IBCABI.Pack("query", queryRes)

	remainingGas, rerr = contract.DeductGas(suppliedGas, cosmosGasUsed)

	return

}

// NewContract returns a new wasmd stateful precompiled contract.
//
//	This contract is used for testing purposes only and should not be used on public chains.
//	The functions of this contract (once implemented), will be used to exercise and test the various aspects of
//	the EVM such as gas usage, argument parsing, events, etc. The specific operations tested under this contract are
//	still to be determined.
func NewContract(wasmdKeeper pcommon.WasmdKeeper, evmKeeper pcommon.EVMKeeper) (contract.StatefulPrecompiledContract, error) {

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
