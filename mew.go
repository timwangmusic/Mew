package main

import (
	"os"
	"strings"

	"astuart.co/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"github.com/weihesdlegend/Mew/config"
	"github.com/weihesdlegend/Mew/transactions"
)

// example input from CLI: mew buy -s 100 -t AAPL
// output: purchased 100 shares of AAPL with market order, and total cost is xxx
func main() {
	log.Info("welcome to use the Mew stock assistant")

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory

	var cfg config.Configurations

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Config read error! %s", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("Unable to decode into struct, %v", err)
	}

	rhClient, err := robinhood.Dial(&robinhood.OAuth{
		Username: cfg.Broker.User,
		Password: cfg.Broker.Password,
	})

	if err != nil {
		log.Fatal("Robinhood auth error %s", err)
	}

	app := cli.NewApp()

	var ticker string
	var shares uint64

	tickerFlag := cli.StringFlag{
		Name:        "ticker",
		Aliases:     []string{"t"},
		Value:       "YANG",
		Required:    false,
		Destination: &ticker,
	}

	sharesFlag := cli.Uint64Flag{
		Name:        "shares",
		Aliases:     []string{"s"},
		Value:       0,
		Required:    false,
		Destination: &shares,
	}

	buyCmd := cli.Command{
		Name:    "buy",
		Aliases: []string{"b"},
		Usage:   "-t MSFT -s 10",
		Flags: []cli.Flag{
			&tickerFlag,
			&sharesFlag,
		},
		Action: func(ctx *cli.Context) error {
			ticker = strings.ToUpper(ticker)
			buyErr := transactions.PlaceMarketOrder(rhClient, ticker, shares, robinhood.Buy)
			if buyErr != nil {
				log.Error(buyErr)
			} else {
				log.Infof("purchased %d shares of %s with %s order", shares, ticker, "market")
			}
			return buyErr
		},
	}

	sellCmd := cli.Command{
		Name:    "sell",
		Usage:   "-t MSFT -s 10",
		Flags: []cli.Flag{
			&tickerFlag,
			&sharesFlag,
		},
		Action: func(ctx *cli.Context) error {
			ticker = strings.ToUpper(ticker)
			sellErr := transactions.PlaceMarketOrder(rhClient, ticker, shares, robinhood.Sell)
			if sellErr != nil {
				log.Error(sellErr)
			} else {
				log.Infof("sold %d shares of %s with %s order", shares, ticker, "market")
			}
			return sellErr
		},
	}

	app.Commands = []*cli.Command{
		&buyCmd,
		&sellCmd,
	}

	appRunErr := app.Run(os.Args)
	if appRunErr != nil {
		log.Fatal(appRunErr)
	}
}
