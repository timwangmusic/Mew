package commands

import (
	"errors"
	"reflect"
	"strings"

	"github.com/coolboy/go-robinhood"
	"github.com/weihesdlegend/Mew/clients"
)

// TODO comment
type MarketSellCommand struct {
	RhClient    clients.Client
	AmountLimit float64
	Ticker      string
	//
	Ins  *robinhood.Instrument
	Opts robinhood.OrderOpts
}

// Readonly
func (base MarketSellCommand) Validate() error {
	if val := reflect.ValueOf(base.RhClient); val.IsZero() {
		return errors.New("RhClient not set")
	}

	if base.AmountLimit <= 0 {
		return errors.New("AmountLimit <= 0")
	}

	if len(base.Ticker) == 0 || len(strings.TrimSpace(base.Ticker)) == 0 {
		return errors.New("ticker cannot be empty")
	}
	return nil
}

// Write, update internal fields
func (base *MarketSellCommand) Prepare() error {
	var err error

	err = base.Validate()
	if err != nil {
		return err
	}

	base.Ins, base.Opts, err = PrepareInsAndOpts(base.Ticker, base.AmountLimit, 100.0, base.RhClient)
	if err != nil {
		return err
	}

	base.Opts.Side = robinhood.Sell
	base.Opts.Type = robinhood.Market

	if err := previewHelper(base.Ticker, base.Opts.Type, base.Opts.Side, base.Opts.Quantity, base.Opts.Price); err != nil {
		return err
	}

	return nil
}

func (base MarketSellCommand) Execute() error {
	return ExecuteOrder(base.Opts, base.Ins, base.RhClient)
}
