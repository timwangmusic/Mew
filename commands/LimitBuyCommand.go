package commands

import (
	"errors"
	"reflect"

	"github.com/weihesdlegend/Mew/clients"

	"github.com/coolboy/go-robinhood"
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

	base.Ins, base.Opts, err = PrepareInsAndOpts(base.Ticker, base.AmountLimit, base.PercentLimit, base.RhClient)
	if err != nil {
		return err
	}

	base.Opts.Side = robinhood.Buy
	base.Opts.Type = robinhood.Limit

	if err = previewHelper(base.Ticker, base.Opts.Type, base.Opts.Side, base.Opts.Quantity, base.Opts.Price); err != nil {
		return err
	}

	return nil
}

func (base LimitBuyCommand) Execute() error {
	return ExecuteOrder(base.Opts, base.Ins, base.RhClient)
}
