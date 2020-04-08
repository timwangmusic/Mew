package commands

import (
	"astuart.co/go-robinhood"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/weihesdlegend/Mew/transactions"
	"strings"
)

// supported commands
var MarketBuyCmd cli.Command
var MarketSellCmd cli.Command
var LimitBuyCmd cli.Command
var LimitSellCmd cli.Command

func InitCommands(rhClient *robinhood.Client) {
	initFlags()

	LimitBuyCmd = cli.Command{
		Name:    "limit",
		Aliases: []string{"lb"},
		Usage:   "-t MSFT -s 10 -l 99.0",
		Flags: []cli.Flag{
			&tickerFlag,
			&sharesFlag,
			&limitFlag,
			&totalValueFlag,
		},
		Action: func(ctx *cli.Context) error {
			if limit > 100.0 {
				err := errors.New("the limit set in limit buy should not exceed 100%")
				return err
			}
			ticker = strings.ToUpper(ticker)
			buyErr, totalVal := transactions.PlaceOrder(rhClient, ticker, shares, robinhood.Buy, robinhood.Limit, totalValue, limit)
			if buyErr != nil {
				log.Error(buyErr)
			} else {
				log.Infof("limit order placed for buying %s with a total value of %.2f", ticker, totalVal)
			}
			return nil
		},
	}

	LimitSellCmd = cli.Command{
		Name:    "limitsell",
		Aliases: []string{"ls"},
		Usage:   "-t MSFT -s 10 -l 101.0",
		Flags: []cli.Flag{
			&tickerFlag,
			&sharesFlag,
			&limitSellFlag,
			&totalValueFlag,
		},
		Action: func(ctx *cli.Context) error {
			if limitSell < 100.0 {
				err := errors.New("the limit set in limit sell should be greater than 100%")
				return err
			}
			ticker = strings.ToUpper(ticker)
			sellErr, totalVal := transactions.PlaceOrder(rhClient, ticker, shares, robinhood.Sell, robinhood.Limit, totalValue, limitSell)
			if sellErr != nil {
				log.Error(sellErr)
			} else {
				log.Infof("limit order placed for selling %s with a total value of %.2f", ticker, totalVal)
			}
			return nil
		},
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
			buyErr, _ := transactions.PlaceOrder(rhClient, ticker, shares, robinhood.Buy, robinhood.Market, 0, 100.0)
			if buyErr != nil {
				log.Error(buyErr)
			} else {
				log.Infof("purchased %d shares of %s with market order", shares, ticker)
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
			sellErr, _ := transactions.PlaceOrder(rhClient, ticker, shares, robinhood.Sell, robinhood.Market, 0, 100.0)
			if sellErr != nil {
				log.Error(sellErr)
			} else {
				log.Infof("sold %d shares of %s with market order", shares, ticker)
			}
			return sellErr
		},
	}
}
