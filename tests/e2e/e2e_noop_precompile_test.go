package e2e_test

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/kava-labs/kava/contracts/contracts/noop_caller"
	"github.com/kava-labs/kava/precompile/contracts/noop"
)

func (suite *IntegrationTestSuite) TestNoopPrecompile_EOAToContractToPrecompile_Success() {
	// setup helper account
	helperAcc := suite.Kava.NewFundedAccount("noop-precompile-eoa->contract->precompile-success-helper-account", sdk.NewCoins(ukava(1e6))) // 1 KAVA

	ethClient, err := ethclient.Dial("http://127.0.0.1:8545")
	suite.Require().NoError(err)

	noopCallerAddr, _, noopCaller, err := noop_caller.DeployNoopCaller(helperAcc.EvmAuth, ethClient)
	suite.Require().NoError(err)

	// wait until contract is deployed
	suite.Eventually(func() bool {
		code, err := ethClient.CodeAt(context.Background(), noopCallerAddr, nil)
		suite.Require().NoError(err)

		return len(code) != 0
	}, 20*time.Second, 1*time.Second)

	code, err := ethClient.CodeAt(context.Background(), noop.ContractAddress, nil)
	suite.Require().NoError(err)
	fmt.Printf("Code: %v\n", code)

	suite.Run("regular type-safe call", func() {
		err := noopCaller.Noop(nil)
		suite.Require().NoError(err)
	})

	suite.Run("evm call & static call opcodes", func() {
		rezBytes, err := noopCaller.NoopStaticCall(nil)
		suite.Require().NoError(err)
		_ = rezBytes
	})
}

func (suite *IntegrationTestSuite) waitForNewBlocks(ethClient *ethclient.Client, n uint64) {
	beginHeight, err := ethClient.BlockNumber(context.Background())
	suite.Require().NoError(err)

	suite.Eventually(func() bool {
		curHeight, err := ethClient.BlockNumber(context.Background())
		suite.Require().NoError(err)

		return curHeight >= beginHeight+n
	}, 10*time.Second, 1*time.Second)
}
