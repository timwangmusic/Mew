package commands

import (
	"errors"
	"github.com/coolboy/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/weihesdlegend/Mew/clients"
	"reflect"
)

type TrailingStopBuyCommand struct {
	Client clients.Client

	Ticker          string
	AmountLimit     float64
	PercentTrailing int // percent trailing the lowest price

	Ins  *robinhood.Instrument
	Opts robinhood.OrderOpts
}

func (cmd TrailingStopBuyCommand) Validate() error {
	if val := reflect.ValueOf(cmd.Client); val.IsZero() {
		return errors.New("broker client is not set")
	}

	if cmd.AmountLimit <= 0 {
		return errors.New("amount limit should be positive")
	}

	if cmd.PercentTrailing <= 0 || cmd.PercentTrailing >= 100.0 {
		return errors.New("percentage in trailing orders should be between 1% and 99%")
	}
	return nil
}

func (cmd *TrailingStopBuyCommand) Prepare() error {
	var err error
	if err = cmd.Validate(); err != nil {
		return err
	}

	var price float64
	cmd.Ins, cmd.Opts, price, err = PrepareInsAndOpts(cmd.Ticker, cmd.AmountLimit, 100.0, cmd.Client)
	if err != nil {
		return err
	}

	cmd.Opts.Side = robinhood.Buy
	if err = TrailingStopOrderHelper(price, &cmd.Opts, cmd.PercentTrailing, cmd.AmountLimit); err != nil {
		return err
	}

	err = trailingStopOrderPreviewHelper(cmd.Ticker, cmd.Opts.Side, cmd.Opts.Quantity, price, cmd.Opts.StopPrice)
	return err
}

func (cmd TrailingStopBuyCommand) Execute() (err error) {
	return ExecuteOrder(cmd.Opts, cmd.Ins, cmd.Client)
}

func TrailingStopBuyCommandCallback(*cli.Context) (err error) {
	client := clients.GetRHClient()

	trailingStopBuyCmd := TrailingStopBuyCommand{
		Client:          client,
		AmountLimit:     totalValue,
		PercentTrailing: percentTrailing,
	}

	var tickers []string
	tickers, err = ParseTicker(ticker)
	for _, ticker := range tickers {
		trailingStopBuyCmd.Ticker = ticker

		err = trailingStopBuyCmd.Prepare()
		if err != nil {
			log.Error(err)
			continue
		}

		err = trailingStopBuyCmd.Execute()
		if err != nil {
			log.Error(err)
			continue
		}
	}
	return
}
