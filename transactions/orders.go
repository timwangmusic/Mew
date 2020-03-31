package transactions

import (
	"astuart.co/go-robinhood"
	"errors"
	log "github.com/sirupsen/logrus"
)

// consider cash account only
func PlaceMarketOrder(client *robinhood.Client, securityName string, quantity uint64, orderSide robinhood.OrderSide) (err error) {
	if quantity == uint64(0) {
		return
	}
	accounts, accountErr := client.GetAccounts()
	if accountErr != nil {
		err = accountErr
		log.Error(err)
		return
	}
	var account = accounts[0]

	quotes, quotesErr := client.GetQuote(securityName)
	if quotesErr != nil {
		err = quotesErr
		return
	} else if len(quotes) == 0 {
		err = errors.New("no quote obtained from provided security name, please check")
		return
	}

	// verify buying power for buy order
	totalValue := quotes[0].AskPrice * float64(quantity)
	if orderSide == robinhood.Buy && totalValue > account.BuyingPower {
		err = errors.New("insufficient buying power")
		log.Errorf("buying power %.2f, need %.2f", account.BuyingPower, totalValue)
		return
	}

	ins, insErr := client.GetInstrumentForSymbol(securityName)
	if insErr != nil {
		err = insErr
		return
	}

	// place order
	_, orderErr := client.Order(ins, robinhood.OrderOpts{Type: robinhood.Market,
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
