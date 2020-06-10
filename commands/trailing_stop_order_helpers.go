package commands

import (
	"errors"
	"github.com/coolboy/go-robinhood"
	"github.com/sirupsen/logrus"
	"github.com/weihesdlegend/Mew/utils"
	"reflect"
)

func TrailingStopOrderHelper(price float64, opts *robinhood.OrderOpts, percentTrailing int, amountLimit float64) error {
	if val := reflect.ValueOf(opts); val.IsNil() {
		return errors.New("invalid order options")
	}

	opts.Type = robinhood.Market

	// time-related settings
	opts.ExtendedHours = false
	opts.TimeInForce = robinhood.GTC // good till cancelled

	// price-related settings
	opts.Trigger = TriggerTypeStop
	opts.TrailingPeg = robinhood.TrailingPeg{
		Type:       TrailingTypePercentage,
		Percentage: percentTrailing,
	}

	if opts.Side == robinhood.Buy {
		opts.StopPrice, _ = utils.Round(price*float64(100+percentTrailing)/100, 0.01)
		opts.Quantity = uint64(amountLimit / opts.StopPrice)
		logrus.Infof("Preparing trailing stop buy order of %d shares with stop price %.2f",
			opts.Quantity, opts.StopPrice)
		opts.Price = opts.StopPrice // this is needed for RH API to execute as trailing stop buy order
	} else if opts.Side == robinhood.Sell {
		opts.StopPrice, _ = utils.Round(price*float64(100-percentTrailing)/100, 0.01)
		opts.Quantity = uint64(amountLimit / opts.StopPrice)
		logrus.Infof("Preparing trailing stop sell order of %d shares with stop price %.2f",
			opts.Quantity, opts.StopPrice)
	}

	return nil
}
