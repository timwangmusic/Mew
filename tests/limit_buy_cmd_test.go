package tests

import (
	"testing"

	"github.com/coolboy/go-robinhood"
	"github.com/stretchr/testify/mock"
	"github.com/weihesdlegend/Mew/commands"
)

var limitBuyCommand = commands.LimitBuyCommand{
	RhClient:     rhClientMocker,
	AmountLimit:  1000.0,
	PercentLimit: 99.0,
	Ticker:       "QQQ",
}

func setupMocker(tickers []string) {
	lastPrice := 100.0
	quotes := []robinhood.Quote{{
		LastTradePrice:              lastPrice,
		LastExtendedHoursTradePrice: lastPrice,
	}}
	for _, ticker := range tickers {
		rhClientMocker.On("GetQuote", ticker).Return(quotes, nil)
	}

	ins := &robinhood.Instrument{}
	rhClientMocker.On("GetInstrument", mock.AnythingOfType("string")).Return(ins, nil)
}

// test limit buy max $1000 worth of stock with limit of 99%
// mock current price at 100.0
// valid case
func TestLimitBuyCommand(t *testing.T) {
	setupMocker([]string{"QQQ"})

	if err := limitBuyCommand.Validate(); err != nil {
		t.Fatal(err)
	}

	if err := limitBuyCommand.Prepare(); err != nil {
		t.Fatal(err)
	}

	lastPrice := 100.0
	expectedLimitPrice := 99.00
	if limitBuyCommand.Opts.Price != limitBuyCommand.PercentLimit*lastPrice/100.0 {
		t.Errorf("expected price to be %.2f, got %.2f", expectedLimitPrice, limitBuyCommand.Opts.Price)
	}

	expectedQuantity := uint64(10)
	if limitBuyCommand.Opts.Quantity != expectedQuantity {
		t.Errorf("expected quantity to be %d, got %d", expectedQuantity, limitBuyCommand.Opts.Quantity)
	}

	rhClientMocker.On("Order", mock.Anything, mock.Anything).Return(&robinhood.OrderOutput{ID: "33521"}, nil)

	if err := limitBuyCommand.Execute(); err != nil {
		t.Fatal(err)
	}

	rhClientMocker.AssertExpectations(t)
}
