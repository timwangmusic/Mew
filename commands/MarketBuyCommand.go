package commands

import (
	"errors"
	"strings"

	"astuart.co/go-robinhood"
	"github.com/weihesdlegend/Mew/clients"
)

// TODO comment
type MarketBuyCommand struct {
	RhClient    clients.Client
	Ticker      string
	AmountLimit float64
	//
	Ins  robinhood.Instrument
	Opts robinhood.OrderOpts
}

// Readonly
func (base MarketBuyCommand) Validate() error {
	// TODO Add validation logic here
	return nil
}

// Write, update internal fields
// TODO should it be stateless?
func (base *MarketBuyCommand) Prepare() error {

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
	limitPrice := quotes[0].Price()
	quantity := uint64(base.AmountLimit / limitPrice)

	base.Opts = robinhood.OrderOpts{
		Type:     robinhood.Market,
		Quantity: quantity,
		Side:     robinhood.Buy,
		Price:    limitPrice,
	}

	return nil
}

func (base MarketBuyCommand) Execute() error {
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
