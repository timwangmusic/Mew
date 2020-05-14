package commands

import (
	"errors"
	"github.com/coolboy/go-robinhood"
	log "github.com/sirupsen/logrus"
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
// TODO should it be stateless?
func (base *MarketBuyCommand) Prepare() error {

	validateErr := base.Validate()
	if validateErr != nil {
		return validateErr
	}

	quotes, quoteErr := base.RhClient.GetQuote(base.Ticker)
	if quoteErr != nil {
		return quoteErr
	}

	if len(quotes) == 0 {
		return errors.New("no quote obtained from provided security name, please check")
	}

	ins, insErr := base.RhClient.GetInstrument(base.Ticker)
	if insErr != nil {
		return insErr
	}

	base.Ins = ins
	price := quotes[0].Price()
	quantity := uint64(base.AmountLimit / price)

	base.Opts = robinhood.OrderOpts{
		Type:          robinhood.Market,
		Quantity:      quantity,
		Side:          robinhood.Buy,
		Price:         price,
		ExtendedHours: true,          // default to allow after hour
		TimeInForce:   robinhood.GFD, // default to GoodForDay
	}

	if err := previewHelper(base.Ticker, base.Opts.Type, base.Opts.Side, base.Opts.Quantity, base.Opts.Price); err != nil {
		return err
	}

	return nil
}

func (base MarketBuyCommand) Execute() error {
	if v := reflect.ValueOf(base.Opts); v.IsZero() {
		return errors.New("please call Prepare()")
	}

	orderRes, orderErr := base.RhClient.MakeOrder(base.Ins, base.Opts)

	if orderErr != nil {
		return orderErr
	}

	log.Infof("Order placed for %s ID %s", base.Ticker, orderRes.ID)

	return nil
}
