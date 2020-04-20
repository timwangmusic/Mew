package tests

import (
	"astuart.co/go-robinhood"
	"flag"
	"github.com/stretchr/testify/mock"
	"github.com/weihesdlegend/Mew/commands"
	"testing"
)

var ticker = flag.String("ticker", "MSFT", "ticker")

func TestPrepare(t *testing.T) {
	mocker := new(RHClientMock)

	quotes := []robinhood.Quote{
		{LastTradePrice: 100.0,
			LastExtendedHoursTradePrice: 105.0},
	}
	mocker.On("GetQuote", mock.Anything).Return(quotes, nil)

	ins := &robinhood.Instrument{}
	mocker.On("GetInstrument", mock.Anything).Return(ins, nil)

	cmd := commands.LimitBuyCommand{
		RhClient: mocker,
	}

	if err := cmd.Prepare(); err != nil {
		t.Error(err)
	}

	mocker.AssertExpectations(t)
}
