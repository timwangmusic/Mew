package commands

import (
	"astuart.co/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/weihesdlegend/Mew/transactions"
	"strings"
)

var ticker string
var shares uint64

var MarketBuyCmd cli.Command
var MarketSellCmd cli.Command

func InitCommands(rhClient *robinhood.Client) {
	tickerFlag := cli.StringFlag{
		Name:        "ticker",
		Aliases:     []string{"t"},
		Value:       "YANG",
		Required:    false,
		Destination: &ticker,
	}

	sharesFlag := cli.Uint64Flag{
		Name:        "shares",
		Aliases:     []string{"s"},
		Value:       0,
		Required:    false,
		Destination: &shares,
	}

	MarketBuyCmd = cli.Command{
		Name:    "buy",
		Aliases: []string{"b"},
		Usage:   "-t MSFT -s 10",
		Flags: []cli.Flag{
			&tickerFlag,
			&sharesFlag,
		},
		Action: func(ctx *cli.Context) error {
			ticker = strings.ToUpper(ticker)
			buyErr := transactions.PlaceMarketOrder(rhClient, ticker, shares, robinhood.Buy)
			if buyErr != nil {
				log.Error(buyErr)
			} else {
				log.Infof("purchased %d shares of %s with %s order", shares, ticker, "market")
			}
			return buyErr
		},
	}

	MarketSellCmd = cli.Command{
		Name:  "sell",
		Usage: "-t MSFT -s 10",
		Flags: []cli.Flag{
			&tickerFlag,
			&sharesFlag,
		},
		Action: func(ctx *cli.Context) error {
			ticker = strings.ToUpper(ticker)
			sellErr := transactions.PlaceMarketOrder(rhClient, ticker, shares, robinhood.Sell)
			if sellErr != nil {
				log.Error(sellErr)
			} else {
				log.Infof("sold %d shares of %s with %s order", shares, ticker, "market")
			}
			return sellErr
		},
	}
}
