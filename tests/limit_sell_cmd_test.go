package tests

import (
	"astuart.co/go-robinhood"
	"github.com/stretchr/testify/mock"
	"github.com/weihesdlegend/Mew/commands"
	"testing"
)

var limitSellCommand = commands.LimitSellCommand{
	RhClient:     rhClientMocker,
	PercentLimit: 101.0,
	AmountLimit:  1010.0,
	Ticker:       "@QQQ_MSFT",
}

// test limit sell $1010 worth of stock with limit of 101%
// mock current price at 100.0
// valid case
func TestLimitSellCommand(t *testing.T) {
	tickers := []string{"QQQ", "MSFT"}
	setupMocker(tickers)

	if err := limitSellCommand.Validate(); err != nil {
		t.Error(err)
	}
	lastPrice := 100.0
	if err := limitSellCommand.Prepare(); err != nil {
		t.Error(err)
	}

	expectedLimitPrice := 101.00
	expectedQuantity := uint64(10)

	for _, ticker := range tickers {
		if limitSellCommand.Opts[ticker].Price != limitSellCommand.PercentLimit*lastPrice/100.0 {
			t.Errorf("expected price to be %.2f, got %.2f", expectedLimitPrice, limitSellCommand.Opts[ticker].Price)
		}
		if limitSellCommand.Opts[ticker].Quantity != expectedQuantity {
			t.Errorf("expected quantity to be %d, got %d", expectedQuantity, limitSellCommand.Opts[ticker].Quantity)
		}
		if limitSellCommand.Opts[ticker].Side != robinhood.Sell {
			t.Errorf("expect side to be sell, got %d", limitSellCommand.Opts[ticker].Side)
		}
		if limitSellCommand.Opts[ticker].Type != robinhood.Limit {
			t.Errorf("expect type to be market, got %d", limitSellCommand.Opts[ticker].Type)
		}
	}

	rhClientMocker.On("MakeOrder", mock.Anything, mock.Anything).Return(&robinhood.OrderOutput{ID: "33522"}, nil)

	if err := limitSellCommand.Execute(); err != nil {
		t.Error(err)
	}

	rhClientMocker.AssertExpectations(t)
}
