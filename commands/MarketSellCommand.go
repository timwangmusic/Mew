package commands

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"

	"astuart.co/go-robinhood"
	"github.com/weihesdlegend/Mew/clients"
)

// TODO comment
type MarketSellCommand struct {
	RhClient    clients.Client
	AmountLimit float64
	Ticker      string
	//
	Ins  map[string]*robinhood.Instrument
	Opts map[string]robinhood.OrderOpts
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
	return nil
}

// Write, update internal fields
func (base *MarketSellCommand) Prepare() error {
	validateErr := base.Validate()
	if validateErr != nil {
		return validateErr
	}

	base.Ins = make(map[string]*robinhood.Instrument)
	base.Opts = make(map[string]robinhood.OrderOpts)

	ticker = strings.ToUpper(base.Ticker)
	tickers := make([]string, 0)
	if strings.HasPrefix(ticker, "@") {
		fields := strings.Split(ticker, "@")
		var batch string
		if len(fields) > 1 {
			batch = fields[1]
		}
		for _, val := range strings.Split(batch, "_") {
			tickers = append(tickers, val)
		}
	} else {
		tickers = append(tickers, ticker)
	}

	if err := base.PrepareInsAndOpts(tickers); err != nil {
		return err
	}

	return nil
}

func (base *MarketSellCommand) PrepareInsAndOpts(tickers []string) (err error) {
	for _, ticker := range tickers {
		ins, insErr := base.RhClient.GetInstrument(ticker)
		if insErr != nil {
			log.Error(insErr)
			err = insErr
			continue
		}
		base.Ins[ticker] = ins

		quotes, quoteErr := base.RhClient.GetQuote(ticker)
		if quoteErr != nil {
			log.Error(quoteErr)
			err = quoteErr
			continue
		}
		if len(quotes) == 0 {
			err = errors.New("no quote obtained from provided security name, please check")
			continue
		}
		price := quotes[0].Price()
		quantity := uint64(base.AmountLimit / price)
		opt := robinhood.OrderOpts{
			Type:          robinhood.Market,
			Quantity:      quantity,
			Side:          robinhood.Sell,
			Price:         price,
			ExtendedHours: true,          // default to allow after hour
			TimeInForce:   robinhood.GFD, // default to GoodForDay
		}
		base.Opts[ticker] = opt
	}
	return
}

func (base MarketSellCommand) Execute() (err error) {
	// place order
	// use ask price in quote to buy or sell
	// time in force defaults to "good for day(gfd)"
	for ticker, ins := range base.Ins {
		if opt, ok := base.Opts[ticker]; ok {
			_, err = base.RhClient.MakeOrder(ins, opt)
		}
	}
	return
}
