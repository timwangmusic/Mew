package commands

import (
	"astuart.co/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/weihesdlegend/Mew/transactions"
)

type LimitBuyCommand struct {
	rhClient *robinhood.Client
	name     string
	aliases  []string
	flags    []cli.Flag
	ticker   string
	limit    float64
	amount   float64
}

func (limitBuy *LimitBuyCommand) Init(rhClient *robinhood.Client, ticker string, limit float64, amount float64) {
	limitBuy.rhClient = rhClient
	limitBuy.name = "limitbuy"
	limitBuy.aliases = []string{"lb"}
	limitBuy.flags = []cli.Flag{
		&tickerFlag,
		&amountFlag,
		&limitBuyFlag,
	}
	limitBuy.ticker = ticker
	limitBuy.limit = limit
	limitBuy.amount = amount
}

func (limitBuy LimitBuyCommand) execute() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		buyErr, totalVal := transactions.PlaceOrder(limitBuy.rhClient, ticker, shares, robinhood.Buy, robinhood.Limit, amount, limit)
		if buyErr != nil {
			log.Error(buyErr)
			return buyErr
		}

		log.Infof("limit order placed for buying %s with a total value of %.2f", ticker, totalVal)
		return nil
	}
}

func (limitBuy LimitBuyCommand) GenCmd() *cli.Command {
	return &cli.Command{
		Name:    limitBuy.name,
		Aliases: limitBuy.aliases,
		Flags:   limitBuy.flags,
		Action:  limitBuy.execute(),
		Usage: "-t MSFT  -l 98.0 -a 2000",
	}
}

// TODO: implement Preview method
func (limitBuy LimitBuyCommand) Preview() Preview {
	return Preview{
		limitBuy.name,
	}
}
