package tests

import (
	"testing"

	"github.com/coolboy/go-robinhood"
	"github.com/stretchr/testify/mock"
	"github.com/weihesdlegend/Mew/commands"
)

var marketSellCommand = commands.MarketSellCommand{
	RhClient:    rhClientMocker,
	AmountLimit: 1000,
	Ticker:      "QQQ",
}

// test market sell $1000 worth of stock
// mock current price at 100.0
// valid case
func TestMarketSellCommand(t *testing.T) {
	tickers := []string{"QQQ"}
	setupMocker(tickers)

	if err := marketSellCommand.Validate(); err != nil {
		t.Error(err)
	}

	if err := marketSellCommand.Prepare(); err != nil {
		t.Error(err)
	}

	expectedQuantity := uint64(10)
	if marketSellCommand.Opts.Quantity != expectedQuantity {
		t.Errorf("expected quantity to be %d, got %d", expectedQuantity, marketSellCommand.Opts.Quantity)
	}
	if marketSellCommand.Opts.Side != robinhood.Sell {
		t.Errorf("expect side to be sell, got %d", marketSellCommand.Opts.Side)
	}
	if marketSellCommand.Opts.Type != robinhood.Market {
		t.Errorf("expect type to be market, got %d", marketSellCommand.Opts.Type)
	}

	rhClientMocker.On("Order", mock.Anything, mock.Anything).Return(&robinhood.OrderOutput{ID: "33524"}, nil)
	if err := marketSellCommand.Execute(); err != nil {
		t.Error(err)
	}

	rhClientMocker.AssertExpectations(t)
}

// test percentage market sell
// valid case
func TestMarketPercentageSell(t *testing.T) {
	marketSellCommand.SellPercent = 50.0
	tickers := []string{"QQQ"}
	setupMocker(tickers)
	setupAdditionalMockerValues(tickers)

	if err := marketSellCommand.Validate(); err != nil {
		t.Error(err)
	}

	if err := marketSellCommand.Prepare(); err != nil {
		t.Error(err)
	}

	expectedQuantity := uint64(5)
	if marketSellCommand.Opts.Quantity != expectedQuantity {
		t.Errorf("expected quantity to be %d, got %d", expectedQuantity, marketSellCommand.Opts.Quantity)
	}
	if marketSellCommand.Opts.Side != robinhood.Sell {
		t.Errorf("expect side to be sell, got %d", marketSellCommand.Opts.Side)
	}
	if marketSellCommand.Opts.Type != robinhood.Market {
		t.Errorf("expect type to be market, got %d", marketSellCommand.Opts.Type)
	}

	rhClientMocker.On("Order", mock.Anything, mock.Anything).Return(&robinhood.OrderOutput{ID: "33522abc"}, nil)
	if err := marketSellCommand.Execute(); err != nil {
		t.Error(err)
	}

	rhClientMocker.AssertExpectations(t)
}
