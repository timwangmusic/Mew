package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"os"

	"github.com/weihesdlegend/Mew/transactions"
	"golang.org/x/oauth2"

	"astuart.co/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"github.com/weihesdlegend/Mew/commands"
	"github.com/weihesdlegend/Mew/config"
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
		log.Error("Unable to decode into struct, %v", err)
		return
	}

	// TODO if cred is empty, please auth
	/*
		ts := &robinhood.OAuth{
			Username: "username",
			Password: "password",
		}

		tk, err := ts.Token()
		// log.Info(tk)

		tkJSON, err := json.Marshal(tk)
		tkJSONb64 := base64.StdEncoding.EncodeToString(tkJSON)
		log.Info(tkJSONb64) // here is your encoded credentials
	*/

	tkJSON, err := base64.StdEncoding.DecodeString(cfg.Broker.EncodedCredentials)
	rawToken := oauth2.Token{}
	json.Unmarshal(tkJSON, &rawToken)

	cts := config.CachedTokenSource{
		RawToken: rawToken,
	}

	cli, err := robinhood.Dial(&cts)

	if err != nil {
		log.Error("Robinhood auth error %s", err)
		os.Exit(1)
	}

	iSPY, err := cli.GetInstrumentForSymbol("SPY")
	log.Info(iSPY)

	// example of a simple market order buy
	buyCmd := flag.NewFlagSet("buy", flag.ExitOnError)
	sellCmd := flag.NewFlagSet("sell", flag.ExitOnError)

	var ticker string
	buyCmd.StringVar(&ticker, "t", "YANG", "stock ticker")
	sellCmd.StringVar(&ticker, "t", "YANG", "stock ticker")

	commands.InitCommands(rhClient)
	app.Commands = []*cli.Command{
		&commands.MarketBuyCmd,
		&commands.MarketSellCmd,
	}

	appRunErr := app.Run(os.Args)
	if appRunErr != nil {
		log.Fatal(appRunErr)
	}
}
