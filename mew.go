package main

import (
	"astuart.co/go-robinhood"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

// example input from CLI: mew buy -s 100 -t AAPL
// output: purchased 100 shares of AAPL with market order, and total cost is 10000
func main() {
	log.Info("welcome to use the mew stock assistant")

	cli, err := robinhood.Dial(&robinhood.OAuth{
		Username: "andrewstuart",
		Password: "mypasswordissecure",
	})

	if err != nil {
		//Oh well
	}

	iSPY, err := cli.GetInstrumentForSymbol("SPY")

	fmt.Print(iSPY)

	// example of a simple market order buy
	buyCmd := flag.NewFlagSet("buy", flag.ExitOnError)

	var ticker string
	buyCmd.StringVar(&ticker, "t", "YANG", "stock ticker")

	var numShares int
	buyCmd.IntVar(&numShares, "s", 0, "number of shares to purchase")

	var orderType string
	buyCmd.StringVar(&orderType, "o", "market", "order type")

	sharePrice := 100 // place holder
	if len(os.Args) > 1 {
		_ = buyCmd.Parse(os.Args[2:])
		ticker = strings.ToUpper(ticker)
		orderType = strings.ToLower(orderType)
	}

	fmt.Printf("purchased %d shares of %s with %s order, and total cost is %d \n", numShares, ticker,
		orderType, numShares*sharePrice)
}
