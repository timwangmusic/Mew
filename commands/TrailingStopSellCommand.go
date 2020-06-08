package commands

import (
	"errors"
	"github.com/coolboy/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/weihesdlegend/Mew/clients"
	"reflect"
)

type TrailingStopSellCommand struct {
	Client clients.Client

	Ticker          string
	AmountLimit     float64
	PercentTrailing int // percent trailing the lowest price

	Ins  *robinhood.Instrument
	Opts robinhood.OrderOpts
}

func (cmd TrailingStopSellCommand) Validate() error {
	if val := reflect.ValueOf(cmd.Client); val.IsZero() {
		return errors.New("broker client is not set")
	}

	if cmd.AmountLimit <= 0 {
		return errors.New("amount limit should be positive")
	}

	if cmd.PercentTrailing <= 0 || cmd.PercentTrailing >= 100.0 {
		return errors.New("percentage in trailing stop orders should be between 1% and 99%")
	}
	return nil
}

func (cmd *TrailingStopSellCommand) Prepare() error {
	var err error
	err = cmd.Validate()
	if err != nil {
		return err
	}

	var price float64
	cmd.Ins, cmd.Opts, price, err = PrepareInsAndOpts(cmd.Ticker, cmd.AmountLimit, 100.0, cmd.Client)
	if err != nil {
		return err
	}

	cmd.Opts.Side = robinhood.Sell

	if err = TrailingStopOrderHelper(price, &cmd.Opts, cmd.PercentTrailing, cmd.AmountLimit); err != nil {
		return err
	}

	err = trailingStopOrderPreviewHelper(cmd.Ticker, cmd.Opts.Side, cmd.Opts.Quantity, price, cmd.Opts.StopPrice)
	return err
}

func (cmd TrailingStopSellCommand) Execute() error {
	return ExecuteOrder(cmd.Opts, cmd.Ins, cmd.Client)
}

func TrailingStopSellCommandCallback(*cli.Context) (err error) {
	client := clients.GetRHClient()

	trailingStopSellCmd := TrailingStopSellCommand{
		Client:          client,
		AmountLimit:     totalValue,
		PercentTrailing: percentTrailing,
	}

	var tickers []string
	tickers, err = ParseTicker(ticker)
	for _, ticker := range tickers {
		trailingStopSellCmd.Ticker = ticker

		err = trailingStopSellCmd.Prepare()
		if err != nil {
			log.Error(err)
			continue
		}

		err = trailingStopSellCmd.Execute()
		if err != nil {
			log.Error(err)
			continue
		}
	}
	return
}
