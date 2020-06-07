package commands

import (
	"errors"
	"github.com/weihesdlegend/Mew/utils"
	"reflect"
	"strings"

	"github.com/coolboy/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/weihesdlegend/Mew/clients"
)

// TODO comment
type LimitSellCommand struct {
	RhClient     clients.Client
	Ticker       string
	PercentLimit float64
	AmountLimit  float64
	PercentSell  float64
	//
	Ins  *robinhood.Instrument
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

	if len(base.Ticker) == 0 || len(strings.TrimSpace(base.Ticker)) == 0 {
		return errors.New("ticker cannot be empty")
	}

	if base.PercentSell < 0 || base.PercentSell > 100.0 {
		return errors.New("sell percent should be greater 0 and no greater than 100.0")
	}
	return nil
}

// Write, update internal fields
func (base *LimitSellCommand) Prepare() error {
	var err error

	err = base.Validate()
	if err != nil {
		return err
	}

	base.Ins, base.Opts, err = ProcessInputsForSell(base.Ticker, base.AmountLimit, base.PercentSell, base.PercentLimit, base.RhClient)
	if err != nil {
		return err
	}

	base.Opts.Side = robinhood.Sell
	base.Opts.Type = robinhood.Limit

	if err := previewHelper(base.Ticker, base.Opts.Type, base.Opts.Side, base.Opts.Quantity, base.Opts.Price); err != nil {
		return err
	}

	return nil
}

func (base LimitSellCommand) Execute() (err error) {
	return ExecuteOrder(base.Opts, base.Ins, base.RhClient)
}

func LimitSellCallback(ctx *cli.Context) (err error) {
	rhClient := clients.GetRHClient()

	tickers, tickerParseErr := ParseTicker(ticker)
	if tickerParseErr != nil {
		err = tickerParseErr
		return
	}

	for _, ticker := range tickers {
		// init
		lsCmd := &LimitSellCommand{
			RhClient:     rhClient,
			Ticker:       ticker,
			AmountLimit:  totalValue,
			PercentLimit: limitSell,
			PercentSell:  percent,
		}
		// preview
		if err = lsCmd.Prepare(); err != nil {
			log.Error("Prepare() for ", ticker, " error : ", err)
			continue
		}

		log.Info(utils.OrderToString(lsCmd.Opts, *lsCmd.Ins, lsCmd.Opts.Price))

		if err = lsCmd.Execute(); err != nil {
			log.Error("Execute() for ", ticker, " error : ", err)
			continue
		}
	}

	return
}
