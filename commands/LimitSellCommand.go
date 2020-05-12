package commands

import (
	"errors"
	"reflect"
	"strings"

	"github.com/coolboy/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/weihesdlegend/Mew/clients"
)

// TODO comment
type LimitSellCommand struct {
	RhClient     clients.Client
	Ticker       string
	PercentLimit float64
	AmountLimit  float64
	//
	Ins  *robinhood.Instrument
	Opts robinhood.OrderOpts
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

	var err error
	base.Ins, base.Opts, err = PrepareInsAndOpts(base.Ticker, base.AmountLimit, base.PercentLimit, base.RhClient)
	if err != nil {
		return err
	}

	base.Opts.Side = robinhood.Sell
	base.Opts.Type = robinhood.Limit
	return nil
}

func (base LimitSellCommand) Execute() (err error) {

	if base.Ins == nil {
		return errors.New("please call Prepare()")
	}

	orderRes, orderErr := base.RhClient.MakeOrder(base.Ins, base.Opts)

	if orderErr != nil {
		return orderErr
	}

	log.Infof("Order placed for %s ID %s", base.Ticker, orderRes.ID)

	return nil
}
