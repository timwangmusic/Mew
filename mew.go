package main

import (
	"flag"
	"fmt"
	"os"
)

// example: ./mew buy --shares 100  --ticker AAPL
// output: purchased 100 shares of AAPL, and total value is 10000
func main() {
	// example of a simple market order buy
	buyCmd := flag.NewFlagSet("buy", flag.ExitOnError)
	//sellCmd := flag.NewFlagSet("sell", flag.ExitOnError)

	ticker := buyCmd.String("ticker", "YANG", "stock ticker")
	numShares := buyCmd.Int("shares", 0, "number of shares to purchase")

	orderType := buyCmd.String("ordertype", "market", "order type")

	if *orderType == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	sharePrice := 100
	_ = buyCmd.Parse(os.Args[2:])

	fmt.Printf("purchased %d shares of %s, and total value is %d \n", *numShares, *ticker, *numShares*sharePrice)
}
