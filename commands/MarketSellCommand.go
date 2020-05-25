package commands

import (
	"errors"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/coolboy/go-robinhood"
	"github.com/weihesdlegend/Mew/clients"
	"github.com/weihesdlegend/Mew/utils"
)

// TODO comment
type MarketSellCommand struct {
	RhClient    clients.Client
	AmountLimit float64
	Ticker      string
	SellPercent  float64
	//
	Ins  *robinhood.Instrument
	Opts robinhood.OrderOpts
}

// Readonly
func (base MarketSellCommand) Validate() error {
	if val := reflect.ValueOf(base.RhClient); val.IsZero() {
		return errors.New("RhClient not set")
	}

	if base.AmountLimit <= 0 {
		return errors.New("AmountLimit <= 0")
	}

	if len(base.Ticker) == 0 || len(strings.TrimSpace(base.Ticker)) == 0 {
		return errors.New("ticker cannot be empty")
	}

	// validate and correct percentage sell flag value
	if base.SellPercent < 0 || base.SellPercent > 100.0 {
		base.SellPercent = 0
	}
	return nil
}

// Write, update internal fields
func (base *MarketSellCommand) Prepare() error {
	var err error

	err = base.Validate()
	if err != nil {
		return err
	}

	base.Ins, base.Opts, err = PrepareInsAndOpts(base.Ticker, base.AmountLimit, 100.0, base.SellPercent, base.RhClient)
	if err != nil {
		return err
	}

	base.Opts.Side = robinhood.Sell
	base.Opts.Type = robinhood.Market

	if err := previewHelper(base.Ticker, base.Opts.Type, base.Opts.Side, base.Opts.Quantity, base.Opts.Price); err != nil {
		return err
	}

	return nil
}

func (base MarketSellCommand) Execute() error {
	return ExecuteOrder(base.Opts, base.Ins, base.RhClient)
}

func MarketSellCallback(ctx *cli.Context) (err error) {
	rhClient := clients.GetRHClient()

	tickers, tickerParseErr := ParseTicker(ticker)
	if tickerParseErr != nil {
		err = tickerParseErr
		return
	}

	for _, ticker := range tickers {
		// init
		msCmd := &MarketSellCommand{
			RhClient:    rhClient,
			Ticker:      ticker,
			AmountLimit: totalValue,
			SellPercent: sellPercent,
		}
		// preview
		if err = msCmd.Prepare(); err != nil {
			continue
		}

		log.Info(utils.OrderToString(msCmd.Opts, *msCmd.Ins))

		if err = msCmd.Execute(); err != nil {
			log.Error("Execute() for ", ticker, " error : ", err)
			continue
		}
	}

	return
}
