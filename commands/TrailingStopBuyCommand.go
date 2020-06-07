package commands

import (
	"errors"
	"github.com/coolboy/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/weihesdlegend/Mew/clients"
	"github.com/weihesdlegend/Mew/utils"
	"reflect"
)

type TrailingStopBuyCommand struct {
	Client clients.Client

	Ticker          string
	AmountLimit     float64
	PercentTrailing int // percent trailing the lowest price

	Ins  *robinhood.Instrument
	Opts robinhood.OrderOpts
}

func (cmd TrailingStopBuyCommand) Validate() error {
	if val := reflect.ValueOf(cmd.Client); val.IsZero() {
		return errors.New("broker client is not set")
	}

	if cmd.AmountLimit <= 0 {
		return errors.New("amount limit should be positive")
	}

	if cmd.PercentTrailing <= 0 || cmd.PercentTrailing >= 100.0 {
		return errors.New("percentage in trailing orders should be between 1% and 99%")
	}
	return nil
}

func (cmd *TrailingStopBuyCommand) Prepare() error {
	var err error
	if err = cmd.Validate(); err != nil {
		return err
	}

	var price float64
	cmd.Ins, cmd.Opts, price, err = PrepareInsAndOpts(cmd.Ticker, cmd.AmountLimit, 100.0, cmd.Client)
	if err != nil {
		return err
	}

	cmd.Opts.Side = robinhood.Buy
	if err = TrailingStopOrderHelper(price, &cmd.Opts, cmd.PercentTrailing, cmd.AmountLimit); err != nil {
		return err
	}

	err = trailingStopOrderPreviewHelper(cmd.Ticker, cmd.Opts.Side, cmd.Opts.Quantity, price, cmd.Opts.StopPrice)
	return err
}

func (cmd TrailingStopBuyCommand) Execute() (err error) {
	return ExecuteOrder(cmd.Opts, cmd.Ins, cmd.Client)
}

func TrailingStopBuyCommandCallback(*cli.Context) (err error) {
	client := clients.GetRHClient()

	trailingStopBuyCmd := TrailingStopBuyCommand{
		Client:          client,
		AmountLimit:     totalValue,
		PercentTrailing: percentTrailing,
	}

	var tickers []string
	tickers, err = ParseTicker(ticker)
	for _, ticker := range tickers {
		trailingStopBuyCmd.Ticker = ticker

		err = trailingStopBuyCmd.Prepare()
		if err != nil {
			log.Error(err)
			continue
		}

		err = trailingStopBuyCmd.Execute()
		if err != nil {
			log.Error(err)
			continue
		}
	}
	return
}

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
		opts.StopPrice, _ = utils.Round(price * float64(100+percentTrailing) / 100, 0.01)
		opts.Quantity = uint64(amountLimit / opts.StopPrice)
		log.Infof("Preparing trailing stop buy order of %d shares with stop price %.2f",
			opts.Quantity, opts.StopPrice)
		opts.Price = opts.StopPrice  // this is needed for RH API to execute as trailing stop buy order
	} else if opts.Side == robinhood.Sell {
		opts.StopPrice, _ = utils.Round(price * float64(100-percentTrailing) / 100, 0.01)
		opts.Quantity = uint64(amountLimit / opts.StopPrice)
		log.Infof("Preparing trailing stop sell order of %d shares with stop price %.2f",
			opts.Quantity, opts.StopPrice)
	}

	return nil
}
