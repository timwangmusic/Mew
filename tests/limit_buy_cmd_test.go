package tests

import (
	"astuart.co/go-robinhood"
	"github.com/stretchr/testify/mock"
	"github.com/weihesdlegend/Mew/commands"
	"testing"
)

var mocker = new(RHClientMock)
var cmd = commands.LimitBuyCommand{
	RhClient: mocker,
	AmountLimit: 1000.0,
	PercentLimit: 99.0,
}

// test limit buy max $1000 worth of stock with limit of 99%
// mock current price at 100.0
// valid case
func TestLimitBuy(t *testing.T) {
	quotes := []robinhood.Quote{
		{
			LastTradePrice: 100.0,
			LastExtendedHoursTradePrice: 100.0,
		},
	}
	mocker.On("GetQuote", mock.AnythingOfType("string")).Return(quotes, nil)

	ins := &robinhood.Instrument{}
	mocker.On("GetInstrument", mock.AnythingOfType("string")).Return(ins, nil)

	if err := cmd.Prepare(); err != nil {
		t.Error(err)
	}

	if cmd.Opts.Price != 99.00 {
		t.Errorf("expected price to be 103.95, got %.2f", cmd.Opts.Price)
	}

	if cmd.Opts.Quantity != 10 {
		t.Errorf("expected quantity to be 10, got %d", cmd.Opts.Quantity)
	}

	mocker.On("MakeOrder", mock.Anything, mock.Anything).Return(&robinhood.OrderOutput{ID: "33521"}, nil)

	cmd.Opts = robinhood.OrderOpts{
		Price: 105,
	}

	if err := cmd.Execute(); err != nil {
		t.Error(err)
	}

	mocker.AssertExpectations(t)
}
