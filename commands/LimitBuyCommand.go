package commands

import (
	"errors"
	"strings"

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
	quotes, quoteErr := base.RhClient.GetQuote(TICK) // TODO make rhClient as interface for testing
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
	quantity := uint64(totalValue / limitPrice)

	base.Opts = robinhood.OrderOpts{
		Type:     robinhood.Limit,
		Quantity: quantity,
		Side:     robinhood.Buy,
		Price:    limitPrice,
	}

	return nil
}

func (base LimitBuyCommand) Execute() error {
	if base.Opts == (robinhood.OrderOpts{}) {
		return errors.New("Please call Prepare()")
	}
	// place order
	// use ask price in quote to buy or sell
	// time in force defaults to "good till canceled(gtc)"
	_, orderErr := base.RhClient.MakeOrder(&base.Ins, base.Opts)

	if orderErr != nil {
		return orderErr
	}

	return nil
}
