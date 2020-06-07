package commands

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strings"

	"github.com/coolboy/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/weihesdlegend/Mew/clients"
	"github.com/weihesdlegend/Mew/utils"
)

const (
	TickerSeparator        = "_"
	TriggerTypeStop        = "stop"
	TrailingTypePercentage = "percentage"
)

// generate order summary for user to confirm
func previewHelper(ticker string, transactionType robinhood.OrderType, side robinhood.OrderSide, quantity uint64, price float64) (err error) {
	// to simplify testing
	if reflect.ValueOf(BufferReader).IsNil() {
		return nil
	}

	fmt.Printf("Please confirm the order details.\n"+
		"Order type: %s %s\t"+
		"Security: %s\t"+
		"Quantity: %d\t"+
		"price: %.2f [y/n]",
		transactionType, side, ticker, quantity, price)

	// wait for user confirmation
	for {
		text, _ := BufferReader.ReadString('\n')
		if strings.Contains(text, "y") {
			break
		} else {
			return errors.New("order is cancelled")
		}
	}
	return nil
}

// parse raw ticker string from user input
func ParseTicker(ticker string) ([]string, error) {
	ticker = strings.ToUpper(ticker)
	tickers := make([]string, 0)
	regex, _ := regexp.Compile(`[A-Z]+`)
	if strings.HasPrefix(ticker, "@") {
		fields := strings.Split(ticker, "@")
		var batch string
		if len(fields) > 1 {
			batch = fields[1]
		}
		for _, val := range strings.Split(batch, TickerSeparator) {
			if regex.MatchString(val) {
				tickers = append(tickers, val)
			}
		}
	} else {
		tickers = append(tickers, ticker)
	}
	if len(tickers) == 0 {
		return tickers, errors.New("ticker parsing error")
	}
	return tickers, nil
}

// process all inputs for sell commands
// the input to this method should capture all the useful flags to the basic sell commands
//
// Params
// ticker: symbol of the security
// amountLimit: maximum value of the security for sell. This value and the order price jointly determine number of shares
// percentLimit: price percentage to apply for limit orders
// percentSell: percentage of the current holdings of the security to sell
func ProcessInputsForSell(ticker string, amountLimit float64, percentSell float64, percentLimit float64, client clients.Client) (Ins *robinhood.Instrument, Opts robinhood.OrderOpts, err error) {
	var price float64
	Ins, Opts, price, err = PrepareInsAndOpts(ticker, amountLimit, percentLimit, client)
	if err != nil {
		return
	}

	Opts.Price = price

	if percentSell > 0 {
		getPositionsCmd, getPositionsCmdErr := GetPositions(client)
		if getPositionsCmdErr != nil {
			err = getPositionsCmdErr
			return
		}

		if _, exist := getPositionsCmd.PositionsMap[ticker]; !exist {
			err = fmt.Errorf("not holding any security for the specified ticker: %s", ticker)
			return
		}
		currentPosition := getPositionsCmd.PositionsMap[ticker]
		// at least sell 1 share in percentage sell mode
		Opts.Quantity = uint64(math.Max(percentSell*currentPosition.Quantity/100, 1.0))
		log.Infof("processing %.2f percent of current holding of %s with a total of %d shares, which is %d shares",
			percentSell, ticker, uint64(currentPosition.Quantity), Opts.Quantity)
	}

	return
}

// make http calls to RH to get instrument data and current security pricing
// generate order options
func PrepareInsAndOpts(ticker string, AmountLimit float64, PercentLimit float64, rhClient clients.Client) (Ins *robinhood.Instrument, Opts robinhood.OrderOpts, finalPrice float64, err error) {
	Ins, insErr := rhClient.GetInstrument(ticker)
	if err = insErr; err != nil {
		return
	}

	quotes, quoteErr := rhClient.GetQuote(ticker)
	if err = quoteErr; err != nil {
		return
	}
	if len(quotes) == 0 {
		err = errors.New("no quote obtained from provided security name, please check")
		return
	}
	price := quotes[0].Price()
	log.Infof("Quote price is %f", price)

	finalPrice, roundErr := utils.Round(price*PercentLimit/100.0, 0.01) // limit to floating point 2 digits
	if err = roundErr; err != nil {
		return
	}

	log.Infof("Updated price is %f", finalPrice)

	quantity := uint64(AmountLimit / finalPrice)

	Opts = robinhood.OrderOpts{
		Quantity:      quantity,
		ExtendedHours: true,          // default to allow after hour
		TimeInForce:   robinhood.GFD, // default to GoodForDay
	}

	return
}

// execute order after obtaining options and instrument details
func ExecuteOrder(opts robinhood.OrderOpts, instrument *robinhood.Instrument, client clients.Client) error {
	if v := reflect.ValueOf(opts); v.IsZero() {
		return errors.New("option is empty, please call Prepare()")
	}

	if reflect.ValueOf(client).IsNil() {
		return errors.New("invalid client")
	}

	orderRes, orderErr := client.Order(instrument, opts)

	if orderErr != nil {
		return orderErr
	}

	log.Infof("Order placed for %s with order ID: %s", instrument.Symbol, orderRes.ID)

	return nil
}
