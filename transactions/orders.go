package transactions

import (
	"astuart.co/go-robinhood"
	"errors"
	log "github.com/sirupsen/logrus"
	"math"
)

func getQuote(client *robinhood.Client, security string) (quotes []robinhood.Quote, err error) {
	quotes, quotesErr := client.GetQuote(security)
	if quotesErr != nil {
		err = quotesErr
		return
	} else if len(quotes) == 0 {
		err = errors.New("no quote obtained from provided security name, please check")
		return
	}
	return
}

func PlaceOrder(client *robinhood.Client, security string, quantity uint64, orderSide robinhood.OrderSide,
	orderType robinhood.OrderType, totalValue float64, limit float64) (err error, totalVal float64) {
	if quantity == 0 && orderType == robinhood.Market {
		return
	}

	quotes, quoteErr := getQuote(client, security)
	if quoteErr != nil {
		err = quoteErr
		return
	}

	ins, insErr := client.GetInstrumentForSymbol(security)
	if insErr != nil {
		err = insErr
		return
	}

	baselinePrice := quotes[0].AskPrice
	totalVal = baselinePrice * float64(quantity)

	orderOptions := robinhood.OrderOpts{
		Type:     orderType,
		Quantity: quantity,
		Side:     orderSide,
		Price:    baselinePrice,
	}

	if orderType == robinhood.Limit {
		limitPrice := round(baselinePrice*limit/100.0, 0.01) // limit to floating point 2 digits
		q := uint64(totalValue / limitPrice)
		orderOptions.Price = limitPrice
		orderOptions.Quantity = q

		log.Debugf("limit price is %.2f, and number of quantity is %d", limitPrice, q)
		totalVal = limitPrice * float64(q)
	}

	// place order
	// use ask price in quote to buy or sell
	// time in force defaults to "good till canceled(gtc)"
	_, orderErr := client.Order(ins, orderOptions)

	if orderErr != nil {
		err = orderErr
		return
	}
	return
}

func round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}
