package tests

import (
	"testing"

	"github.com/coolboy/go-robinhood"
	"github.com/stretchr/testify/mock"
	"github.com/weihesdlegend/Mew/commands"
)

var limitSellCommand = commands.LimitSellCommand{
	RhClient:     rhClientMocker,
	PercentLimit: 101.0,
	AmountLimit:  1010.0,
	Ticker:       "QQQ",
}

// test limit sell $1010 worth of stock with limit of 101%
// mock current price at 100.0
// valid case
func TestLimitSellCommand(t *testing.T) {
	tickers := []string{"QQQ"}
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

	if limitSellCommand.Opts.Price != limitSellCommand.PercentLimit*lastPrice/100.0 {
		t.Errorf("expected price to be %.2f, got %.2f", expectedLimitPrice, limitSellCommand.Opts.Price)
	}
	if limitSellCommand.Opts.Quantity != expectedQuantity {
		t.Errorf("expected quantity to be %d, got %d", expectedQuantity, limitSellCommand.Opts.Quantity)
	}
	if limitSellCommand.Opts.Side != robinhood.Sell {
		t.Errorf("expect side to be sell, got %d", limitSellCommand.Opts.Side)
	}
	if limitSellCommand.Opts.Type != robinhood.Limit {
		t.Errorf("expect type to be market, got %d", limitSellCommand.Opts.Type)
	}

	rhClientMocker.On("MakeOrder", mock.Anything, mock.Anything).Return(&robinhood.OrderOutput{ID: "33522"}, nil)

	if err := limitSellCommand.Execute(); err != nil {
		t.Error(err)
	}

	rhClientMocker.AssertExpectations(t)
}
