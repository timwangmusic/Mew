package tests

import (
	"astuart.co/go-robinhood"
	"fmt"
	"github.com/stretchr/testify/mock"
)

type RHClientMock struct {
	mock.Mock
}

func (c *RHClientMock) GetQuote(ticker string) (quotes []robinhood.Quote, err error) {
	fmt.Printf("method called with %s", ticker)
	args := c.Called(ticker)
	return args.Get(0).([]robinhood.Quote), args.Error(1)
}

func (c *RHClientMock) GetInstrument(ticker string) (ins *robinhood.Instrument, err error) {
	fmt.Printf("method called with %s", ticker)
	args := c.Called(ticker)
	return args.Get(0).(*robinhood.Instrument), args.Error(1)
}

func (c *RHClientMock) MakeOrder(ins *robinhood.Instrument, opts robinhood.OrderOpts) (orderOutput *robinhood.OrderOutput, err error) {
	return
}
