package commands

import (
	"errors"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/weihesdlegend/Mew/clients"

	"astuart.co/go-robinhood"
)

// TODO comment
type LimitBuyCommand struct {
	RhClient     clients.Client
	Ticker       string
	PercentLimit float64
	AmountLimit  float64
	//
	Ins  robinhood.Instrument
	Opts robinhood.OrderOpts
}

// Readonly
func (base LimitBuyCommand) Validate() error {
	// TODO Add validation logic here
	return nil
}

// Write, update internal fields
// TODO should it be stateless?
func (base *LimitBuyCommand) Prepare() error {

	validateErr := base.Validate()
	if validateErr != nil {
		return validateErr
	}

	TICK := strings.ToUpper(ticker)
	quotes, quoteErr := base.RhClient.GetQuote(TICK)
	if quoteErr != nil {
		return quoteErr
	}

	if len(quotes) == 0 {
		return errors.New("no quote obtained from provided security name, please check")
	}

	ins, insErr := base.RhClient.GetInstrument(TICK)
	if insErr != nil {
		return insErr
	}
	base.Ins = *ins

	baselinePrice := quotes[0].Price()
	limitPrice := round(baselinePrice*base.PercentLimit/100.0, 0.01) // limit to floating point 2 digits
	quantity := uint64(base.AmountLimit / limitPrice)

	base.Opts = robinhood.OrderOpts{
		Type:     robinhood.Limit,
		Quantity: quantity,
		Side:     robinhood.Buy,
		Price:    limitPrice,
	}

	log.Infof("About to place %d shares of %s buy order at price %.2f", base.Opts.Quantity, base.Opts.Type, base.Opts.Price)
	return nil
}

func (base LimitBuyCommand) Execute() error {
	if v := reflect.ValueOf(base.Opts); v.IsZero() {
		return errors.New("please call Prepare()")
	}
	// place order
	// use ask price in quote to buy or sell
	// time in force defaults to "good till canceled(gtc)"
	orderRes, orderErr := base.RhClient.MakeOrder(&base.Ins, base.Opts)

	if orderErr != nil {
		return orderErr
	}

	log.Infof("Order placed with order ID %s", orderRes.ID)

	return nil
}
