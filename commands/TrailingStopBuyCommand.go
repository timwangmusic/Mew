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

	cmd.Ins, cmd.Opts, err = PrepareInsAndOpts(cmd.Ticker, cmd.AmountLimit, 100.0, cmd.Client)
	if err != nil {
		return err
	}

	cmd.Opts.Side = robinhood.Buy
	cmd.Opts.Type = robinhood.Market

	// price-related settings
	cmd.Opts.Trigger = TriggerTypeStop
	cmd.Opts.TrailingPeg = robinhood.TrailingPeg{
		Type:       TrailingTypePercentage,
		Percentage: cmd.PercentTrailing,
	}
	cmd.Opts.StopPrice, _ = utils.Round(cmd.Opts.Price*float64(100+cmd.PercentTrailing)/100.0, 0.01)
	cmd.Opts.Quantity = uint64(cmd.AmountLimit / cmd.Opts.StopPrice)
	log.Infof("updated quantity with stop price of %.2f is %d shares", cmd.Opts.StopPrice, cmd.Opts.Quantity)

	// time-related settings
	cmd.Opts.ExtendedHours = false
	cmd.Opts.TimeInForce = robinhood.GTC // good till cancelled

	err = previewHelper(ticker, cmd.Opts.Type, cmd.Opts.Side, cmd.Opts.Quantity, cmd.Opts.Price)
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
		log.Info(utils.OrderToString(trailingStopBuyCmd.Opts, *trailingStopBuyCmd.Ins))

		err = trailingStopBuyCmd.Execute()
		if err != nil {
			log.Error(err)
			continue
		}
	}
	return
}
