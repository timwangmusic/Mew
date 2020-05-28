package tests

import (
	"github.com/coolboy/go-robinhood"
	"github.com/stretchr/testify/mock"
)

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

func setupAdditionalMockerValues(tickers []string) {
	var positions []robinhood.Position
	p := robinhood.Position{
		AverageBuyPrice: 100.0,
		Instrument:      "",
		Quantity:        10,
	}
	for _, ticker := range tickers {
		rhClientMocker.On("GetInstrumentByURL", mock.AnythingOfType("string")).
			Return(&robinhood.Instrument{Symbol: ticker}, nil)
		positions = append(positions, p)
	}

	rhClientMocker.On("GetPositions").Return(positions, nil)
}
