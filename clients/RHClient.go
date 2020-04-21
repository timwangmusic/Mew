package clients

import (
	"astuart.co/go-robinhood"
	"golang.org/x/oauth2"
)

type Client interface {
	GetQuote(ticker string) ([]robinhood.Quote, error)
	GetInstrument(ticker string) (*robinhood.Instrument, error)
	MakeOrder(*robinhood.Instrument, robinhood.OrderOpts) (*robinhood.OrderOutput, error)
}

type RHClient struct {
	Client *robinhood.Client
}

func (c *RHClient) Init(token oauth2.TokenSource) (err error) {
	c.Client, err = robinhood.Dial(token)
	return
}

func (c *RHClient) GetQuote(ticker string) (quotes []robinhood.Quote, err error) {
	quotes, err = c.Client.GetQuote(ticker)
	return
}

func (c *RHClient) GetInstrument(ticker string) (ins *robinhood.Instrument, err error) {
	ins, err = c.Client.GetInstrumentForSymbol(ticker)
	return
}

func (c *RHClient) MakeOrder(ins *robinhood.Instrument, opts robinhood.OrderOpts) (orderOutput *robinhood.OrderOutput, err error) {
	orderOutput, err = c.Client.Order(ins, opts)
	return
}
