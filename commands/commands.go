package commands

import (
	"astuart.co/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/weihesdlegend/Mew/clients"
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
			lbCmd := &LimitBuyCommand{
				RhClient:     &clients.RHClient{Client: rhClient},
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
		},
		Action: func(ctx *cli.Context) error {
			// init
			lsCmd := &LimitSellCommand{
				RhClient:     &clients.RHClient{Client: rhClient},
				Ticker:       ticker,
				PercentLimit: limitSell,
				AmountLimit:  totalValue,
			}
			// TODO show preview here
			lsCmd.Prepare()
			// Exec
			sellErr := lsCmd.Execute()
			if sellErr != nil {
				log.Error(sellErr)
				return sellErr
			}

			return nil
		},
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
		Action: func(ctx *cli.Context) error {
			// init
			mbCmd := &MarketBuyCommand{
				rhClient:    rhClient,
				Ticker:      ticker,
				AmountLimit: totalValue,
			}
			// TODO show preview here
			mbCmd.Prepare()
			// Exec
			buyErr := mbCmd.Execute()
			if buyErr != nil {
				log.Error(buyErr)
				return buyErr
			}

			return nil
		},
	}

	MarketSellCmd = cli.Command{
		Name:    "sell",
		Aliases: []string{"s"},
		Usage:   "-t QQQ -v 250",
		Flags: []cli.Flag{
			&tickerFlag,
			// TODO &sharesFlag,
			&totalValueFlag,
		},
		Action: func(ctx *cli.Context) error {
			// init
			msCmd := &MarketSellCommand{
				rhClient:    rhClient,
				Ticker:      ticker,
				AmountLimit: totalValue,
			}
			// TODO show preview here
			msCmd.Prepare()
			// Exec
			err := msCmd.Execute()
			if err != nil {
				log.Error(err)
				return err
			}

			return nil
		},
	}
}
