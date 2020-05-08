package commands

import (
	"errors"
	"github.com/weihesdlegend/Mew/utils"
	"reflect"
	"strings"

	"astuart.co/go-robinhood"
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
	Ins  robinhood.Instrument
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

	return nil
}

// Write, update internal fields
func (base *LimitSellCommand) Prepare() error {

	validateErr := base.Validate()
	if validateErr != nil {
		return validateErr
	}

	TICK := strings.ToUpper(base.Ticker)
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
	log.Infof("quoted price is %f", baselinePrice)
	limitPrice, roundErr := utils.Round(baselinePrice*base.PercentLimit/100.0, 0.01) // limit to floating point 2 digits
	if roundErr != nil {
		return roundErr
	}
	log.Infof("limit price is %f", limitPrice)
	quantity := uint64(base.AmountLimit / limitPrice)

	base.Opts = robinhood.OrderOpts{
		Type:     robinhood.Limit,
		Quantity: quantity,
		Side:     robinhood.Sell,
		Price:    limitPrice,

		ExtendedHours: true,          // default to allow after hour
		TimeInForce:   robinhood.GFD, // default to GoodForDay
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
	orderRes, orderErr := base.RhClient.MakeOrder(&base.Ins, base.Opts)

	if orderErr != nil {
		return orderErr
	}

	log.Infof("Order placed with order ID %s", orderRes.ID)

	return nil
}
