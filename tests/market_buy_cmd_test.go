package tests

import (
	"testing"

	"github.com/coolboy/go-robinhood"
	"github.com/stretchr/testify/mock"
	"github.com/weihesdlegend/Mew/commands"
)

var marketBuyCommand = commands.MarketBuyCommand{
	RhClient:    rhClientMocker,
	AmountLimit: 1000.0,
	Ticker:      "QQQ",
}

// test market buy max $1000 worth of stock
// mock current price at 100.0
// valid case
func TestMarketBuyCommand(t *testing.T) {
	setupMocker([]string{"QQQ"})

	if err := marketBuyCommand.Validate(); err != nil {
		t.Error(err)
	}
	if err := marketBuyCommand.Prepare(); err != nil {
		t.Error(err)
	}

	expectedQuantity := uint64(10)
	if marketBuyCommand.Opts.Quantity != expectedQuantity {
		t.Errorf("expected quantity to be %d, got %d", expectedQuantity, marketBuyCommand.Opts.Quantity)
	}

	rhClientMocker.On("Order", mock.Anything, mock.Anything).Return(&robinhood.OrderOutput{ID: "33523"}, nil)
	if err := marketBuyCommand.Execute(); err != nil {
		t.Error(err)
	}

	rhClientMocker.AssertExpectations(t)
}
