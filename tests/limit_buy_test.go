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
	PercentLimit: 101.0,
}

func TestLimitBuy(t *testing.T) {
	quotes := []robinhood.Quote{
		{LastTradePrice: 100.0,
			LastExtendedHoursTradePrice: 105.0},
	}
	mocker.On("GetQuote", mock.AnythingOfType("string")).Return(quotes, nil)

	ins := &robinhood.Instrument{}
	mocker.On("GetInstrument", mock.AnythingOfType("string")).Return(ins, nil)

	if err := cmd.Prepare(); err != nil {
		t.Error(err)
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
