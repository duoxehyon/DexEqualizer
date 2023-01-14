package main

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tarik0/DexEqualizer/abis"
	"github.com/tarik0/DexEqualizer/logger"
	"math/big"
	"os"
	"strings"
)

func main() {
	rawAbi := abis.FlashloanExecutorV2MetaData.ABI
	execAbi, err := abi.JSON(strings.NewReader(rawAbi))
	if err != nil {
		logger.Log.WithError(err).Fatalln("Unable to load flashloan executor abi.")
	}

	if len(os.Args) < 2 {
		logger.Log.Fatalln("Usage: decode_params.go 0x...")
	}
	paramsBytes, err := hex.DecodeString(os.Args[1][10:])
	if err != nil {
		logger.Log.WithError(err).Fatalln("Unable to decode parameters.")
	}

	params, err := execAbi.Methods["executeFlashloan"].Inputs.Unpack(paramsBytes)
	if err != nil {
		logger.Log.WithError(err).Fatalln("Unable to unpack parameters.")
	}

	tx_param := params[0].(struct {
		Pairs                 []common.Address   "json:\"Pairs\""
		Reserves              [][]*big.Int       "json:\"Reserves\""
		PairTokens            [][]common.Address "json:\"PairTokens\""
		Path                  []common.Address   "json:\"Path\""
		AmountsOut            []*big.Int         "json:\"AmountsOut\""
		BorrowFee             *big.Int           "json:\"BorrowFee\""
		RevertOnReserveChange bool               "json:\"RevertOnReserveChange\""
	})

	logger.Log.Infoln("Pairs:")
	for i, pair := range tx_param.Pairs {
		logger.Log.Infoln(fmt.Sprintf("[%d/%d] %s", i, len(tx_param.Pairs), pair.String()))
	}
	logger.Log.Infoln("")

	logger.Log.Infoln("Path")
	for i, addr := range tx_param.Path {
		logger.Log.Infoln(fmt.Sprintf("[%d/%d] %s", i, len(tx_param.Pairs), addr.String()))
	}

	logger.Log.Infoln("")

	logger.Log.Infoln("Amounts:")
	for i, amount := range tx_param.AmountsOut {
		logger.Log.Infoln(fmt.Sprintf("[%d/%d] %s", i, len(tx_param.AmountsOut), amount.String()))
	}

	logger.Log.Infoln("")

	logger.Log.Infoln("Pair Tokens:")
	for i, tokens := range tx_param.PairTokens {
		logger.Log.Infoln(fmt.Sprintf("[%d/%d] %s, %s", i, len(tx_param.PairTokens), tokens[0].String(), tokens[1].String()))
	}
}
