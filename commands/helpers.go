package commands

import (
	"astuart.co/go-robinhood"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/weihesdlegend/Mew/clients"
	"github.com/weihesdlegend/Mew/utils"
	"strings"
)

func ParseTicker(ticker string) []string {
	ticker = strings.ToUpper(ticker)
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
	return tickers
}

func PrepareInsAndOpts(tickers []string, AmountLimit float64, PercentLimit float64, rhClient clients.Client) (Ins map[string]*robinhood.Instrument,
	Opts map[string]*robinhood.OrderOpts, err error) {
	Ins = make(map[string]*robinhood.Instrument)
	Opts = make(map[string]*robinhood.OrderOpts)

	for _, ticker := range tickers {
		ins, insErr := rhClient.GetInstrument(ticker)
		if insErr != nil {
			logrus.Error(insErr)
			err = insErr
			continue
		}
		Ins[ticker] = ins

		quotes, quoteErr := rhClient.GetQuote(ticker)
		if quoteErr != nil {
			logrus.Error(quoteErr)
			err = quoteErr
			continue
		}
		if len(quotes) == 0 {
			err = errors.New("no quote obtained from provided security name, please check")
			continue
		}
		price := quotes[0].Price()
		price, roundErr := utils.Round(price*PercentLimit/100.0, 0.01) // limit to floating point 2 digits
		if roundErr != nil {
			err = roundErr
			continue
		}
		quantity := uint64(AmountLimit / price)
		opt := robinhood.OrderOpts{
			Quantity:      quantity,
			Price:         price,
			ExtendedHours: true,          // default to allow after hour
			TimeInForce:   robinhood.GFD, // default to GoodForDay
		}
		Opts[ticker] = &opt
	}
	return
}

