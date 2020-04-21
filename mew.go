package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/weihesdlegend/Mew/commands"
)

// example input from CLI: mew buy -s 100 -t AAPL
// output: purchased 100 shares of AAPL with market order, and total cost is xxx
func main() {
	log.Info("welcome to use the Mew stock assistant")

	// TODO if cred is empty, please auth

	/*
		ts := &robinhood.OAuth{
			Username: "user",
			Password: "pwd",
		}

		tk, err := ts.Token()

		tkJSON, err := json.Marshal(tk)
		tkJSONb64 := base64.StdEncoding.EncodeToString(tkJSON)
		log.Info(tkJSONb64) // here is your encoded credentials
	*/

	commands.InitCommands()

	app := cli.NewApp()
	app.Commands = []*cli.Command{
		&commands.MarketBuyCmd,
		&commands.MarketSellCmd,
		&commands.LimitBuyCmd,
		&commands.LimitSellCmd,
	}

	appRunErr := app.Run(os.Args)
	if appRunErr != nil {
		log.Fatal(appRunErr)
	}
}
