package transactions

import (
	"astuart.co/go-robinhood"
	"errors"
)

// consider cash account only
func PlaceMarketOrder(client *robinhood.Client, securityName string, quantity uint64, orderSide robinhood.OrderSide) (err error) {
	if quantity == uint64(0) {
		return
	}

	quotes, quotesErr := client.GetQuote(securityName)
	if quotesErr != nil {
		err = quotesErr
		return
	} else if len(quotes) == 0 {
		err = errors.New("no quote obtained from provided security name, please check")
		return
	}

	ins, insErr := client.GetInstrumentForSymbol(securityName)
	if insErr != nil {
		err = insErr
		return
	}

	// place order
	// use ask price in quote to buy or sell
	// time in force defaults to "good till canceled(gtc)"
	_, orderErr := client.Order(ins, robinhood.OrderOpts{
		Type:     robinhood.Market,
		Quantity: quantity,
		Side:     orderSide,
		Price:    quotes[0].AskPrice,
	})

	if orderErr != nil {
		err = orderErr
		return
	}
	return
}
