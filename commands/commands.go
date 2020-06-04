package commands

import (
	"bufio"
	"os"

	"github.com/urfave/cli/v2"
)

// supported commands
var MarketBuyCmd cli.Command
var MarketSellCmd cli.Command
var LimitBuyCmd cli.Command
var LimitSellCmd cli.Command
var AuthCmd cli.Command
var TrailingStopBuyCmd cli.Command

var BufferReader *bufio.Reader

func InitCommands() {
	BufferReader = bufio.NewReader(os.Stdin)
	AuthCmd = cli.Command{
		Name:    "authenticate",
		Aliases: []string{"auth"},
		Usage:   "-u usr -p pwd -m 1111",
		Flags: []cli.Flag{
			&userFlag,
			&passwordFlag,
			&mfaFlag,
		},
		Action: AuthCallback,
	}

	TrailingStopBuyCmd = cli.Command{
		Name: "trailing_stop_buy",
		Aliases: []string{"tsb"},
		Usage:   "-t MSFT -pt 10 -v 2000",
		Flags: []cli.Flag{
			&tickerFlag,
			&percentTrailingFlag,
			&totalValueFlag,
		},
		Action: TrailingStopBuyCommandCallback,
	}

	LimitBuyCmd = cli.Command{
		Name:    "limitbuy",
		Aliases: []string{"lb"},
		Usage:   "-t MSFT -l 99.0 -v 2000",
		Flags: []cli.Flag{
			&tickerFlag,
			&sharesFlag,
			&limitBuyFlag,
			&totalValueFlag,
		},
		Action: LimitBuyCallback,
	}

	// TODO -ls needs to be updated
	LimitSellCmd = cli.Command{
		Name:    "limitsell",
		Aliases: []string{"ls"},
		Usage:   "-t MSFT -ls 101.0 -v 2000",
		Flags: []cli.Flag{
			&tickerFlag,
			&sharesFlag,
			&limitSellFlag,
			&totalValueFlag,
			&percentToSellFlag,
		},
		Action: LimitSellCallback,
	}

	MarketBuyCmd = cli.Command{
		Name:    "buy",
		Aliases: []string{"b"},
		Usage:   "-t MSFT -v 200",
		Flags: []cli.Flag{
			&tickerFlag,
			// TODO &sharesFlag,
			&totalValueFlag,
		},
		Action: MarketBuyCallback,
	}

	MarketSellCmd = cli.Command{
		Name:    "sell",
		Aliases: []string{"s"},
		Usage:   "-t QQQ -v 250",
		Flags: []cli.Flag{
			&tickerFlag,
			// TODO &sharesFlag,
			&totalValueFlag,
			&percentToSellFlag,
		},
		Action: MarketSellCallback,
	}
}
