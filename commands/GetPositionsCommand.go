package commands

import (
	"errors"
	"github.com/coolboy/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/weihesdlegend/Mew/clients"
	"reflect"
	"sync"
)

type GetPositionsCommand struct {
	RhClient  clients.Client
	Positions []Position
	PositionsMap map[string]Position
}

func (cmd *GetPositionsCommand) Validate() error {
	if val := reflect.ValueOf(cmd.RhClient); val.IsNil() {
		return errors.New("client is not set")
	}
	return nil
}

func (cmd *GetPositionsCommand) Prepare() error {
	return cmd.Validate()
}

func (cmd *GetPositionsCommand) Execute() error {
	positions, err := cmd.RhClient.GetPositions()
	if err != nil {
		return err
	}

	cmd.Positions = make([]Position, len(positions))
	cmd.PositionsMap = make(map[string]Position)

	wg := sync.WaitGroup{}
	wg.Add(len(positions))
	// concurrently get details for each position
	for idx, p := range positions {
		go getPosition(&wg, cmd.Positions, idx, p, cmd.RhClient)
	}
	wg.Wait()

	for _, p := range cmd.Positions {
		if len(p.Ticker) > 0 && p.Quantity > 0 {
			cmd.PositionsMap[p.Ticker] = p
			log.Debugf("ticker: %s, quote price: %.2f, quantity: %.2f", p.Ticker, p.QuotePrice, p.Quantity)
		}
	}
	return nil
}

// private method
func getPosition(wg *sync.WaitGroup, positions []Position, idx int, p robinhood.Position, client clients.Client) {
	defer wg.Done()

	var newPosition Position
	newPosition.AverageBuyPrice = p.AverageBuyPrice
	newPosition.Quantity = p.Quantity

	ins, insErr := client.GetInstrumentByURL(p.Instrument)
	if insErr != nil {
		log.Error(insErr)
		return
	}

	newPosition.Ticker = ins.Symbol
	quotes, quotesErr := client.GetQuote(ins.Symbol)
	if quotesErr != nil {
		log.Error(quotesErr)
		return
	}

	newPosition.QuotePrice = quotes[0].Price()
	newPosition.Valid = true
	positions[idx] = newPosition
	return
}

func GetPositions(cli clients.Client) (*GetPositionsCommand, error) {
	getPositionsCmd := GetPositionsCommand{
		RhClient: cli,
	}

	var err error

	err = getPositionsCmd.Prepare()
	if err != nil {
		return &getPositionsCmd, err
	}

	err = getPositionsCmd.Execute()
	return &getPositionsCmd, err
}
