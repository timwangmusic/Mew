package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/weihesdlegend/Mew/commands"
)

var (
	// Version The app version
	Version = "v0.1.3.1"
	// BuildDate The app build date in yyyy-mm-dd
	BuildDate = "2020-05-31"
)

// example input from CLI: mew buy -s 100 -t AAPL
// output: purchased 100 shares of AAPL with market order, and total cost is xxx
func main() {
	log.Info("welcome to use the Mew stock assistant")

	commands.InitCommands()

	app := &cli.App{
		Name:    "Mew",
		Version: Version + "_build_" + BuildDate,
	}

	app.Commands = []*cli.Command{
		&commands.MarketBuyCmd,
		&commands.MarketSellCmd,
		&commands.LimitBuyCmd,
		&commands.LimitSellCmd,
		&commands.AuthCmd, // Create Auth file
	}

	appRunErr := app.Run(os.Args)
	if appRunErr != nil {
		log.Fatal(appRunErr)
	}
}
