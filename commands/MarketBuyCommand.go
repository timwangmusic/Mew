package commands

import (
	"errors"
	"github.com/coolboy/go-robinhood"
	"github.com/weihesdlegend/Mew/clients"
	"reflect"
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

	base.Ins, base.Opts, err = PrepareInsAndOpts(base.Ticker, base.AmountLimit, 100.0, base.RhClient)
	if err != nil {
		return err
	}

	base.Opts.Side = robinhood.Buy
	base.Opts.Type = robinhood.Market

	if err = previewHelper(base.Ticker, base.Opts.Type, base.Opts.Side, base.Opts.Quantity, base.Opts.Price); err != nil {
		return err
	}

	return nil
}

func (base MarketBuyCommand) Execute() error {
	return ExecuteOrder(base.Opts, base.Ins, base.RhClient)
}
