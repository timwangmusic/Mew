package main

import (
	"astuart.co/go-robinhood"
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/weihesdlegend/Mew/transactions"
	"os"
	"strings"
)

// example input from CLI: mew buy -s 100 -t AAPL
// output: purchased 100 shares of AAPL with market order, and total cost is 10000
func main() {
	log.Info("welcome to use the Mew stock assistant")

	cli, err := robinhood.Dial(&robinhood.OAuth{
		Username: "andrewstuart",
		Password: "mypasswordissecure",
	})

	if err != nil {
		log.Error(err)
	}

	// example of a simple market order buy
	buyCmd := flag.NewFlagSet("buy", flag.ExitOnError)

	var ticker string
	buyCmd.StringVar(&ticker, "t", "YANG", "stock ticker")

	var numShares uint64
	buyCmd.Uint64Var(&numShares, "n", 0, "number of shares to purchase")

	var orderType string
	buyCmd.StringVar(&orderType, "o", "market", "order type")

	if len(os.Args) > 1 {
		_ = buyCmd.Parse(os.Args[2:])
		ticker = strings.ToUpper(ticker)
		orderType = strings.ToLower(orderType)
	}

	if !(robinhood.IsRegularTradingTime() || robinhood.IsExtendedTradingTime()) {
		log.Info("out of regular trading or extended trading time, order will be fulfilled later")
	}

	buyErr := transactions.PlaceMarketOrder(cli, ticker, numShares, robinhood.Buy)
	if buyErr != nil {
		log.Error(buyErr)
	} else {
		log.Infof("purchased %d shares of %s with %s order", numShares, ticker, orderType)
	}
}
