package commands

import (
	"errors"
	"reflect"
	"strings"

	"astuart.co/go-robinhood"
	"github.com/weihesdlegend/Mew/clients"
)

// TODO comment
type LimitSellCommand struct {
	RhClient     clients.Client
	Ticker       string
	PercentLimit float64
	AmountLimit  float64
	//
	Ins  map[string]*robinhood.Instrument
	Opts map[string]*robinhood.OrderOpts
}

// Readonly
func (base LimitSellCommand) Validate() error {
	if val := reflect.ValueOf(base.RhClient); val.IsZero() {
		return errors.New("RhClient not set")
	}

	if base.AmountLimit <= 0 {
		return errors.New("AmountLimit <= 0")
	}

	if base.PercentLimit <= 0 {
		return errors.New("PercentLimit <= 0")
	}

	if len(base.Ticker) == 0 || len(strings.TrimSpace(base.Ticker)) == 0 {
		return errors.New("ticker cannot be empty")
	}

	return nil
}

// Write, update internal fields
func (base *LimitSellCommand) Prepare() error {
	validateErr := base.Validate()
	if validateErr != nil {
		return validateErr
	}

	base.Ins = make(map[string]*robinhood.Instrument)
	base.Opts = make(map[string]*robinhood.OrderOpts)

	tickers := ParseTicker(base.Ticker)

	var err error
	base.Ins, base.Opts, err = PrepareInsAndOpts(tickers, base.AmountLimit, base.PercentLimit, base.RhClient)
	if err != nil {
		return err
	}

	for _, opt := range base.Opts {
		opt.Side = robinhood.Sell
		opt.Type = robinhood.Limit
	}

	return nil
}

func (base LimitSellCommand) Execute() (err error) {
	for ticker, ins := range base.Ins {
		if opt, ok := base.Opts[ticker]; ok {
			_, err = base.RhClient.MakeOrder(ins, *opt)
		}
	}
	return nil
}
