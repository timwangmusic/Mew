package commands

import (
	"errors"
	"github.com/weihesdlegend/Mew/utils"
	"reflect"

	"github.com/coolboy/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/weihesdlegend/Mew/clients"
)

// TODO comment
type MarketBuyCommand struct {
	RhClient    clients.Client
	Ticker      string
	AmountLimit float64
	//
	Ins  *robinhood.Instrument
	Opts robinhood.OrderOpts
}

// Readonly
func (base MarketBuyCommand) Validate() error {
	if val := reflect.ValueOf(base.RhClient); val.IsZero() {
		return errors.New("RhClient not set")
	}

	if base.AmountLimit <= 0 {
		return errors.New("AmountLimit <= 0")
	}

	return nil
}

// Write, update internal fields
func (base *MarketBuyCommand) Prepare() error {
	var err error

	err = base.Validate()
	if err != nil {
		return err
	}

	var price float64
	base.Ins, base.Opts, price, err = PrepareInsAndOpts(base.Ticker, base.AmountLimit, 100.0, base.RhClient)
	if err != nil {
		return err
	}

	base.Opts.Side = robinhood.Buy
	base.Opts.Type = robinhood.Market
	base.Opts.Price = price

	if err = previewHelper(base.Ticker, base.Opts.Type, base.Opts.Side, base.Opts.Quantity, base.Opts.Price); err != nil {
		return err
	}

	return nil
}

func (base MarketBuyCommand) Execute() error {
	return ExecuteOrder(base.Opts, base.Ins, base.RhClient)
}

func MarketBuyCallback(ctx *cli.Context) (err error) {
	rhClient := clients.GetRHClient()
	// init
	mbCmd := &MarketBuyCommand{
		RhClient:    rhClient,
		Ticker:      ticker,
		AmountLimit: totalValue,
	}

	var tickers []string
	tickers, err = ParseTicker(ticker)
	if err != nil {
		return
	}

	for _, ticker := range tickers {
		mbCmd.Ticker = ticker

		// prepare and preview
		err = mbCmd.Prepare()
		if err != nil {
			log.Error(err)
			continue
		}
		log.Info(utils.OrderToString(mbCmd.Opts, *mbCmd.Ins, mbCmd.Opts.Price))

		// execution
		err = mbCmd.Execute()
		if err != nil {
			log.Error(err)
			continue
		}
	}
	return
}
