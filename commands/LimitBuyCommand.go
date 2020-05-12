package commands

import (
	"errors"
	"reflect"
	"strings"

	"github.com/weihesdlegend/Mew/utils"

	log "github.com/sirupsen/logrus"
	"github.com/weihesdlegend/Mew/clients"

	"github.com/coolboy/go-robinhood"
)

// TODO comment
type LimitBuyCommand struct {
	RhClient     clients.Client
	Ticker       string
	PercentLimit float64 // 99.95 stand for 99.95%
	AmountLimit  float64

	//
	Ins  robinhood.Instrument
	Opts robinhood.OrderOpts
}

// Readonly
func (base LimitBuyCommand) Validate() error {
	if val := reflect.ValueOf(base.RhClient); val.IsZero() {
		return errors.New("RhClient not set")
	}

	if base.AmountLimit <= 0 {
		return errors.New("AmountLimit <= 0")
	}

	if base.PercentLimit <= 0 || base.PercentLimit >= 150 {
		return errors.New("PercentLimit <= 0 || >= 150")
	}

	return nil
}

// Write, update internal fields
// TODO should it be stateless?
func (base *LimitBuyCommand) Prepare() error {

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
		Side:     robinhood.Buy,
		Price:    limitPrice,
		// TODO, make it into config? or env_var
		ExtendedHours: true,          // default to allow after hour
		TimeInForce:   robinhood.GFD, // default to GoodForDay
	}

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
