package json_test

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/precompile/contracts/json"
	"github.com/kava-labs/kava/precompile/registry"
	"github.com/stretchr/testify/require"
	"github.com/tharsis/ethermint/x/evm/statedb"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	tmtypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

func MockAddressPair() (sdk.AccAddress, common.Address) {
	return PrivateKeyToAddresses(MockPrivateKey())
}

func MockPrivateKey() cryptotypes.PrivKey {
	// Generate a new Sei private key
	entropySeed, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropySeed)
	algo := hd.Secp256k1
	derivedPriv, _ := algo.Derive()(mnemonic, "", "")
	return algo.Generate()(derivedPriv)
}

func PrivateKeyToAddresses(privKey cryptotypes.PrivKey) (sdk.AccAddress, common.Address) {
	// Encode the private key to hex (i.e. what wallets do behind the scene when users reveal private keys)
	testPrivHex := hex.EncodeToString(privKey.Bytes())

	// Sign an Ethereum transaction with the hex private key
	key, _ := crypto.HexToECDSA(testPrivHex)
	msg := crypto.Keccak256([]byte("foo"))
	sig, _ := crypto.Sign(msg, key)

	// Recover the public keys from the Ethereum signature
	recoveredPub, _ := crypto.Ecrecover(msg, sig)
	pubKey, _ := crypto.UnmarshalPubkey(recoveredPub)

	return sdk.AccAddress(privKey.PubKey().Address()), crypto.PubkeyToAddress(*pubKey)
}

func TestExtractAsBytes(t *testing.T) {
	tApp := app.NewTestApp()
	ctx := tApp.NewContext(true, tmtypes.Header{Height: 1, ChainID: "kava-test", Time: time.Now().UTC()})

	evm := vm.EVM{
		StateDB: statedb.New(ctx, tApp.GetEvmKeeper(), statedb.NewEmptyTxConfig(common.BytesToHash(ctx.HeaderHash().Bytes()))),
	}
	p, err := json.NewContract()
	require.Nil(t, err)
	method := json.ABI.Methods[json.ExtractAsBytesMethod]
	suppliedGas := uint64(10_000_000)

	for _, test := range []struct {
		body           []byte
		expectedOutput []byte
	}{
		{
			[]byte("{\"key\":1}"),
			[]byte("1"),
		}, {
			[]byte("{\"key\":\"1\"}"),
			[]byte("1"),
		}, {
			[]byte("{\"key\":[1,2,3]}"),
			[]byte("[1,2,3]"),
		}, {
			[]byte("{\"key\":{\"nested\":1}}"),
			[]byte("{\"nested\":1}"),
		},
	} {
		args, err := method.Inputs.Pack(test.body, "key")
		require.Nil(t, err)
		res, _, err := p.Run(&evm, registry.JsonContractAddress, registry.JsonContractAddress,
			append(method.ID, args...),
			suppliedGas,
			false,
			nil,
		)
		require.Nil(t, err)
		output, err := method.Outputs.Unpack(res)
		require.Nil(t, err)
		require.Equal(t, 1, len(output))
		require.Equal(t, output[0].([]byte), test.expectedOutput)
	}
}

func TestExtractAsBytesList(t *testing.T) {
	tApp := app.NewTestApp()
	ctx := tApp.NewContext(true, tmtypes.Header{Height: 1, ChainID: "kava-test", Time: time.Now().UTC()})

	evm := vm.EVM{
		StateDB: statedb.New(ctx, tApp.GetEvmKeeper(), statedb.NewEmptyTxConfig(common.BytesToHash(ctx.HeaderHash().Bytes()))),
	}
	p, err := json.NewContract()
	require.Nil(t, err)
	method := json.ABI.Methods[json.ExtractAsBytesListMethod]
	suppliedGas := uint64(10_000_000)

	for _, test := range []struct {
		body           []byte
		expectedOutput [][]byte
	}{
		{
			[]byte("{\"key\":[],\"key2\":1}"),
			[][]byte{},
		}, {
			[]byte("{\"key\":[1,2,3],\"key2\":1}"),
			[][]byte{[]byte("1"), []byte("2"), []byte("3")},
		}, {
			[]byte("{\"key\":[\"1\", \"2\"],\"key2\":1}"),
			[][]byte{[]byte("\"1\""), []byte("\"2\"")},
		}, {
			[]byte("{\"key\":[{\"nested\":1}, {\"nested\":2}],\"key2\":1}"),
			[][]byte{[]byte("{\"nested\":1}"), []byte("{\"nested\":2}")},
		},
	} {
		args, err := method.Inputs.Pack(test.body, "key")
		require.Nil(t, err)
		res, _, err := p.Run(&evm, registry.JsonContractAddress, registry.JsonContractAddress,
			append(method.ID, args...),
			suppliedGas,
			false,
			nil,
		)
		require.Nil(t, err)
		output, err := method.Outputs.Unpack(res)
		require.Nil(t, err)
		require.Equal(t, 1, len(output))
		require.Equal(t, output[0].([][]byte), test.expectedOutput)
	}
}

func TestExtractAsUint256(t *testing.T) {
	tApp := app.NewTestApp()
	ctx := tApp.NewContext(true, tmtypes.Header{Height: 1, ChainID: "kava-test", Time: time.Now().UTC()})

	evm := vm.EVM{
		StateDB: statedb.New(ctx, tApp.GetEvmKeeper(), statedb.NewEmptyTxConfig(common.BytesToHash(ctx.HeaderHash().Bytes()))),
	}
	p, err := json.NewContract()
	require.Nil(t, err)
	method := json.ABI.Methods[json.ExtractAsUint256Method]
	suppliedGas := uint64(10_000_000)
	n := new(big.Int)

	n.SetString("12345678901234567890", 10)
	for _, test := range []struct {
		body           []byte
		expectedOutput *big.Int
	}{
		{
			[]byte("{\"key\":\"12345678901234567890\"}"),
			n,
		}, {
			[]byte("{\"key\":\"0\"}"),
			big.NewInt(0),
		},
	} {
		args, err := method.Inputs.Pack(test.body, "key")
		require.Nil(t, err)
		res, _, err := p.Run(&evm, registry.JsonContractAddress, registry.JsonContractAddress,
			append(method.ID, args...),
			suppliedGas,
			false,
			nil,
		)
		require.Nil(t, err)
		output, err := method.Outputs.Unpack(res)
		require.Nil(t, err)
		require.Equal(t, 1, len(output))
		require.Equal(t, 0, output[0].(*big.Int).Cmp(test.expectedOutput))
	}
}
