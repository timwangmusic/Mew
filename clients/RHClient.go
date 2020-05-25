package clients

import (
	"github.com/coolboy/go-robinhood"
	"golang.org/x/oauth2"
)

type Client interface {
	GetQuote(ticker string) ([]robinhood.Quote, error)
	GetInstrument(ticker string) (*robinhood.Instrument, error)

	// We should rename to Order() to align to internal robinhood calls
	Order(*robinhood.Instrument, robinhood.OrderOpts) (*robinhood.OrderOutput, error)

	GetPositions() ([]robinhood.Position, error)

	GetInstrumentByURL(string) (*robinhood.Instrument, error)
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

func (c *RHClient) Order(ins *robinhood.Instrument, opts robinhood.OrderOpts) (orderOutput *robinhood.OrderOutput, err error) {
	orderOutput, err = c.Client.Order(ins, opts)
	return
}

func (c *RHClient) GetPositions() (positions []robinhood.Position, err error) {
	positions, err = c.Client.GetPositions()
	return
}

func (c *RHClient) GetInstrumentByURL(url string) (ins *robinhood.Instrument, err error) {
	ins, err = c.Client.GetInstrument(url)
	return
}
