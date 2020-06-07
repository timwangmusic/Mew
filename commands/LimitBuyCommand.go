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
type LimitBuyCommand struct {
	RhClient     clients.Client
	Ticker       string
	PercentLimit float64 // 99.95 stand for 99.95%
	AmountLimit  float64

	//
	Ins  *robinhood.Instrument
	Opts robinhood.OrderOpts
}

// Readonly
func (base LimitBuyCommand) Validate() error {
	if val := reflect.ValueOf(base.RhClient); val.IsZero() {
		return errors.New("RhClient not set")
	}

	if base.AmountLimit <= 0 {
		return errors.New("AmountLimit <= 0")
	}

	if base.PercentLimit <= 0 || base.PercentLimit >= 150 {
		return errors.New("PercentLimit <= 0 || >= 150")
	}

	return nil
}

// Write, update internal fields
func (base *LimitBuyCommand) Prepare() error {
	var err error

	err = base.Validate()
	if err != nil {
		return err
	}

	var price float64
	base.Ins, base.Opts, price, err = PrepareInsAndOpts(base.Ticker, base.AmountLimit, base.PercentLimit, base.RhClient)
	if err != nil {
		return err
	}

	base.Opts.Side = robinhood.Buy
	base.Opts.Type = robinhood.Limit
	base.Opts.Price = price

	if err = previewHelper(base.Ticker, base.Opts.Type, base.Opts.Side, base.Opts.Quantity, base.Opts.Price); err != nil {
		return err
	}

	return nil
}

func (base LimitBuyCommand) Execute() error {
	return ExecuteOrder(base.Opts, base.Ins, base.RhClient)
}

// TODO Should we consolidate LimitBuyCallback && LimitSellCallback?
// LimitBuyCallback to refactor the code
func LimitBuyCallback(ctx *cli.Context) (err error) {
	rhClient := clients.GetRHClient()

	var tickers []string
	tickers, err = ParseTicker(ticker)
	if err != nil {
		return
	}

	// init
	lbCmd := &LimitBuyCommand{
		RhClient:     rhClient,
		Ticker:       ticker,
		PercentLimit: limit,
		AmountLimit:  totalValue,
	}

	for _, ticker := range tickers {
		lbCmd.Ticker = ticker
		// prepare and preview
		err = lbCmd.Prepare()
		if err != nil {
			log.Error(err)
			continue
		}
		log.Info(utils.OrderToString(lbCmd.Opts, *lbCmd.Ins))

		// execution
		err = lbCmd.Execute()
		if err != nil {
			log.Error("Execute() for ", ticker, " error : ", err)
			continue
		}
	}

	return
}
