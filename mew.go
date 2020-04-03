package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"astuart.co/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/weihesdlegend/Mew/config"
)

// example input from CLI: mew buy -s 100 -t AAPL
// output: purchased 100 shares of AAPL with market order, and total cost is 10000
func main() {
	log.Info("welcome to use the mew stock assistant")

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory

	var cfg config.Configurations

	if err := viper.ReadInConfig(); err != nil {
		log.Error("Config read error! %s", err)
		return
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Info("Unable to decode into struct, %v", err)
	}

	cli, err := robinhood.Dial(&robinhood.OAuth{
		Username: cfg.Broker.User,
		Password: cfg.Broker.Password,
	})

	if err != nil {
		log.Error("Robinhood auth error %s", err)
	}

	iSPY, err := cli.GetInstrumentForSymbol("SPY")

	fmt.Print(iSPY)

	// example of a simple market order buy
	buyCmd := flag.NewFlagSet("buy", flag.ExitOnError)

	var ticker string
	buyCmd.StringVar(&ticker, "t", "YANG", "stock ticker")

	var numShares int
	buyCmd.IntVar(&numShares, "s", 0, "number of shares to purchase")

	var orderType string
	buyCmd.StringVar(&orderType, "o", "market", "order type")

	sharePrice := 100 // place holder
	if len(os.Args) > 1 {
		_ = buyCmd.Parse(os.Args[2:])
		ticker = strings.ToUpper(ticker)
		orderType = strings.ToLower(orderType)
	}

	fmt.Printf("purchased %d shares of %s with %s order, and total cost is %d \n", numShares, ticker,
		orderType, numShares*sharePrice)
}
