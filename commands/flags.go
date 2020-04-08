package commands

import "github.com/urfave/cli/v2"

// flag destination variables
var ticker string
var shares uint64
var limit float64
var limitSell float64
var totalValue float64

// supported flags
var tickerFlag cli.StringFlag
var sharesFlag cli.Uint64Flag
var limitFlag cli.Float64Flag
var limitSellFlag cli.Float64Flag
var totalValueFlag cli.Float64Flag

func initFlags() {
	tickerFlag = cli.StringFlag{
		Name:        "ticker",
		Aliases:     []string{"t"},
		Value:       "YANG",
		Required:    false,
		Destination: &ticker,
	}

	sharesFlag = cli.Uint64Flag{
		Name:        "shares",
		Aliases:     []string{"s"},
		Value:       0,
		Required:    false,
		Destination: &shares,
	}

	limitFlag = cli.Float64Flag{
		Name:        "limitforbuy",
		Aliases:     []string{"l"},
		Required:    false,
		Value:       99.0,
		Destination: &limit,
	}

	limitSellFlag = cli.Float64Flag{
		Name:        "limitforsell",
		Aliases:     []string{"ls"},
		Required:    false,
		Value:       101.0,
		Destination: &limitSell,
	}

	totalValueFlag = cli.Float64Flag{
		Name:        "value",
		Aliases:     []string{"v"},
		Required:    false,
		Value:       500.0,
		Destination: &totalValue,
	}
}
