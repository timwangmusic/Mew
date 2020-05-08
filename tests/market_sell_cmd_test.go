package tests

import (
	"astuart.co/go-robinhood"
	"github.com/stretchr/testify/mock"
	"github.com/weihesdlegend/Mew/commands"
	"testing"
)

var marketSellCommand = commands.MarketSellCommand{
	RhClient:    rhClientMocker,
	AmountLimit: 1000,
	Ticker:      "@QQQ_MSFT",
}

// test market sell $1000 worth of stock
// mock current price at 100.0
// valid case
func TestMarketSellCommand(t *testing.T) {
	tickers := []string{"QQQ", "MSFT"}
	setupMocker(tickers)

	if err := marketSellCommand.Validate(); err != nil {
		t.Error(err)
	}

	if err := marketSellCommand.Prepare(); err != nil {
		t.Error(err)
	}

	expectedQuantity := uint64(10)
	for _, ticker := range tickers {
		if marketSellCommand.Opts[ticker].Quantity != expectedQuantity {
			t.Errorf("expected quantity to be %d, got %d", expectedQuantity, marketSellCommand.Opts[ticker].Quantity)
		}
		if marketSellCommand.Opts[ticker].Side != robinhood.Sell {
			t.Errorf("expect side to be sell, got %d", marketSellCommand.Opts[ticker].Side)
		}
		if marketSellCommand.Opts[ticker].Type != robinhood.Market {
			t.Errorf("expect type to be market, got %d", marketSellCommand.Opts[ticker].Type)
		}
	}

	rhClientMocker.On("MakeOrder", mock.Anything, mock.Anything).Return(&robinhood.OrderOutput{ID: "33524"}, nil)
	if err := marketSellCommand.Execute(); err != nil {
		t.Error(err)
	}

	rhClientMocker.AssertExpectations(t)
}
