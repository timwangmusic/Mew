package commands

import "github.com/urfave/cli/v2"

// flag destination variables
var ticker string
var shares uint64
var limit  float64
var totalValue float64

// supported flags
var tickerFlag cli.StringFlag
var sharesFlag cli.Uint64Flag
var limitFlag cli.Float64Flag
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
		Name:        "limit",
		Aliases:     []string{"l"},
		Required:    false,
		Value:       100.0,
		Destination: &limit,
	}

	totalValueFlag = cli.Float64Flag{
		Name:        "value",
		Aliases:     []string{"v"},
		Required:    false,
		Value:       500.0,
		Destination: &totalValue,
	}
}
