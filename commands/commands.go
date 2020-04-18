package commands

import (
	"astuart.co/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/weihesdlegend/Mew/transactions"
)

// supported commands
var MarketBuyCmd cli.Command
var MarketSellCmd cli.Command
var LimitBuyCmd cli.Command
var LimitSellCmd cli.Command

func InitCommands(rhClient *robinhood.Client) {
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
		Action: func(ctx *cli.Context) error {
			// init
			var lbCmd CommandBase
			lbCmd = &LimitBuyCommand{
				rhClient:     rhClient,
				Ticker:       ticker,
				PercentLimit: limit,
				AmountLimit:  totalValue,
			}
			// TODO show preview here
			lbCmd.Prepare()
			// Exec
			buyErr := lbCmd.Execute()
			if buyErr != nil {
				log.Error(buyErr)
				return buyErr
			}

			return nil
		},
	}

	LimitSellCmd = cli.Command{
		Name:    "limitsell",
		Aliases: []string{"ls"},
		Usage:   "-t MSFT  -l 101.0 -v 2000",
		Flags: []cli.Flag{
			&tickerFlag,
			&sharesFlag,
			&limitSellFlag,
			&totalValueFlag,
		},
		Action: func(ctx *cli.Context) error {
			sellErr, totalVal := transactions.PlaceOrder(rhClient, ticker, shares, robinhood.Sell, robinhood.Limit, totalValue, limitSell)
			if sellErr != nil {
				log.Error(sellErr)
				return sellErr
			}

			log.Infof("limit order placed for selling %s with a total value of %.2f", ticker, totalVal)
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
		Name:    "sell",
		Aliases: []string{"s"},
		Usage:   "-t MSFT -s 10",
		Flags: []cli.Flag{
			&tickerFlag,
			&sharesFlag,
		},
		Action: func(ctx *cli.Context) error {
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
