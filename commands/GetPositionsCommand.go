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
	RhClient     clients.Client
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
	rawPositions, err := cmd.RhClient.GetPositions()
	if err != nil {
		return err
	}

	positions := make([]Position, len(rawPositions))
	cmd.PositionsMap = make(map[string]Position)

	wg := sync.WaitGroup{}
	wg.Add(len(rawPositions))
	// concurrently get details for each position
	for idx, position := range rawPositions {
		go getPosition(&wg, positions, idx, position, cmd.RhClient)
	}
	wg.Wait()

	for _, position := range positions {
		if len(position.Ticker) > 0 && position.Quantity > 0 {
			cmd.PositionsMap[position.Ticker] = position
			log.Debugf("ticker: %s, quote price: %.2f, quantity: %.2f", position.Ticker, position.QuotePrice, position.Quantity)
		}
	}
	return nil
}

// private method
func getPosition(wg *sync.WaitGroup, positions []Position, idx int, position robinhood.Position, client clients.Client) {
	defer wg.Done()

	var newPosition Position
	newPosition.AverageBuyPrice = position.AverageBuyPrice
	newPosition.Quantity = position.Quantity

	ins, insErr := client.GetInstrumentByURL(position.Instrument)
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
