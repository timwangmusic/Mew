package tests

import (
	"astuart.co/go-robinhood"
	"github.com/stretchr/testify/mock"
)

type RHClientMock struct {
	mock.Mock
}

func (c *RHClientMock) GetQuote(ticker string) (quotes []robinhood.Quote, err error) {
	args := c.Called(ticker)
	return args.Get(0).([]robinhood.Quote), args.Error(1)
}

func (c *RHClientMock) GetInstrument(ticker string) (ins *robinhood.Instrument, err error) {
	args := c.Called(ticker)
	return args.Get(0).(*robinhood.Instrument), args.Error(1)
}

func (c *RHClientMock) MakeOrder(ins *robinhood.Instrument, opts robinhood.OrderOpts) (orderOutput *robinhood.OrderOutput, err error) {
	args := c.Called(ins, opts)
	return args.Get(0).(*robinhood.OrderOutput), args.Error(1)
}

var rhClientMocker = new(RHClientMock)
