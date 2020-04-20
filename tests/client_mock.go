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
	return
}

func (c *RHClientMock) GetInstrument(ticker string) (ins *robinhood.Instrument, err error) {
	fmt.Printf("method called with %s", ticker)
	return
}

func (c *RHClientMock) MakeOrder(ins *robinhood.Instrument, opts robinhood.OrderOpts) (orderOutput *robinhood.OrderOutput, err error) {
	return
}
