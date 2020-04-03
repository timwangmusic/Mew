package main

import (
	"flag"
	"github.com/weihesdlegend/Mew/transactions"
	"os"
	"strings"

	"astuart.co/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/weihesdlegend/Mew/config"
)

// example input from CLI: mew buy -s 100 -t AAPL
// output: purchased 100 shares of AAPL with market order, and total cost is xxx
func main() {
	log.Info("welcome to use the Mew stock assistant")

	if len(os.Args) < 2 {
		log.Info("no transaction type specified")
		os.Exit(1)
	}

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
    os.Exit(1)
	}

	// example of a simple market order buy
	buyCmd := flag.NewFlagSet("buy", flag.ExitOnError)
	sellCmd := flag.NewFlagSet("sell", flag.ExitOnError)

	var ticker string
	buyCmd.StringVar(&ticker, "t", "YANG", "stock ticker")
	sellCmd.StringVar(&ticker, "t", "YANG", "stock ticker")

	var numShares uint64
	buyCmd.Uint64Var(&numShares, "n", 0, "number of shares to purchase")
	sellCmd.Uint64Var(&numShares, "n", 0, "number of shares to purchase")

	var orderType string
	buyCmd.StringVar(&orderType, "o", "market", "order type")
	sellCmd.StringVar(&orderType, "o", "market", "order type")

	switch os.Args[1] {
	case "buy":
		_ = buyCmd.Parse(os.Args[2:])
	case "sell":
		_ = sellCmd.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	ticker = strings.ToUpper(ticker)
	orderType = strings.ToLower(orderType)

	if buyCmd.Parsed(){
		buyErr := transactions.PlaceMarketOrder(cli, ticker, numShares, robinhood.Buy)
		if buyErr != nil {
			log.Error(buyErr)
		} else {
			log.Infof("purchased %d shares of %s with %s order", numShares, ticker, orderType)
		}
	} else if sellCmd.Parsed() {
		sellErr := transactions.PlaceMarketOrder(cli, ticker, numShares, robinhood.Sell)
		if sellErr != nil {
			log.Error(sellErr)
		} else {
			log.Infof("sold %d shares of %s with %s order", numShares, ticker, orderType)
		}
	}
}
