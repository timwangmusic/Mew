package commands

import (
	"errors"
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
	Ins  robinhood.Instrument
	Opts robinhood.OrderOpts
}

// Readonly
func (base LimitSellCommand) Validate() error {
	// TODO Add validation logic here
	return nil
}

// Write, update internal fields
// TODO should it be stateless?
func (base *LimitSellCommand) Prepare() error {

	validateErr := base.Validate()
	if validateErr != nil {
		return validateErr
	}

	TICK := strings.ToUpper(ticker)
	quotes, quoteErr := base.rhClient.GetQuote(TICK) // TODO make rhClient as interface for testing
	if quoteErr != nil {
		return quoteErr
	}

	if len(quotes) == 0 {
		return errors.New("no quote obtained from provided security name, please check")
	}

	ins, insErr := base.rhClient.GetInstrumentForSymbol(TICK)
	if insErr != nil {
		return insErr
	}
	base.Ins = *ins

	baselinePrice := quotes[0].Price()
	limitPrice := Util{}.round(baselinePrice*base.PercentLimit/100.0, 0.01) // limit to floating point 2 digits
	quantity := uint64(totalValue / limitPrice)

	base.Opts = robinhood.OrderOpts{
		Type:     robinhood.Limit,
		Quantity: quantity,
		Side:     robinhood.Sell,
		Price:    limitPrice,
	}

	return nil
}

func (base LimitSellCommand) Execute() error {
	if base.Opts == (robinhood.OrderOpts{}) {
		return errors.New("Please call Prepare()")
	}
	// place order
	// use ask price in quote to buy or sell
	// time in force defaults to "good till canceled(gtc)"
	_, orderErr := base.rhClient.Order(&base.Ins, base.Opts)

	if orderErr != nil {
		return orderErr
	}

	return nil
}
