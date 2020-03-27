package main

import (
	"flag"
	"fmt"
)

func main() {
	// example of a simple market order
	ticker := flag.String("ticker", "YANG", "stock ticker")
	numShares := flag.Int("shares", 0, "number of shares to purchase")
	sharePrice := 100
	flag.Parse()

	fmt.Printf("purchased %d shares of %s, and total value is %d \n", *numShares, *ticker, *numShares*sharePrice)
}
