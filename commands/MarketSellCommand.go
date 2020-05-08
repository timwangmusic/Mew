package commands

import (
	"errors"
	"reflect"
	"strings"

	"astuart.co/go-robinhood"
	"github.com/weihesdlegend/Mew/clients"
)

// TODO comment
type MarketSellCommand struct {
	RhClient    clients.Client
	AmountLimit float64
	Ticker      string
	//
	Ins  map[string]*robinhood.Instrument
	Opts map[string]*robinhood.OrderOpts
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
	validateErr := base.Validate()
	if validateErr != nil {
		return validateErr
	}

	base.Ins = make(map[string]*robinhood.Instrument)
	base.Opts = make(map[string]*robinhood.OrderOpts)

	tickers := ParseTicker(base.Ticker)

	var err error
	base.Ins, base.Opts, err = PrepareInsAndOpts(tickers, base.AmountLimit, 100.0, base.RhClient)
	if err != nil {
		return err
	}

	for _, opt := range base.Opts {
		opt.Side = robinhood.Sell
		opt.Type = robinhood.Market
	}
	return nil
}

func (base MarketSellCommand) Execute() (err error) {
	for ticker, ins := range base.Ins {
		if opt, ok := base.Opts[ticker]; ok {
			_, err = base.RhClient.MakeOrder(ins, *opt)
		}
	}
	return
}
