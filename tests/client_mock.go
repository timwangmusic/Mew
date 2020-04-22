package tests

import (
	"astuart.co/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type RHClientMock struct {
	mock.Mock
}

func (c *RHClientMock) GetQuote(ticker string) (quotes []robinhood.Quote, err error) {
	log.Infof("GetQuote called with %s", ticker)
	args := c.Called(ticker)
	return args.Get(0).([]robinhood.Quote), args.Error(1)
}

func (c *RHClientMock) GetInstrument(ticker string) (ins *robinhood.Instrument, err error) {
	log.Infof("GetInstrument called with %s", ticker)
	args := c.Called(ticker)
	return args.Get(0).(*robinhood.Instrument), args.Error(1)
}

func (c *RHClientMock) MakeOrder(ins *robinhood.Instrument, opts robinhood.OrderOpts) (orderOutput *robinhood.OrderOutput, err error) {
	log.Infof("MakeOrder called")
	args := c.Called(ins, opts)
	return args.Get(0).(*robinhood.OrderOutput), args.Error(1)
}
