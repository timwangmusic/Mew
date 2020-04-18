package main

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"astuart.co/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"github.com/weihesdlegend/Mew/commands"
	"github.com/weihesdlegend/Mew/config"
	"golang.org/x/oauth2"
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
			Username: "user",
			Password: "pwd",
		}

		tk, err := ts.Token()

		tkJSON, err := json.Marshal(tk)
		tkJSONb64 := base64.StdEncoding.EncodeToString(tkJSON)
		log.Info(tkJSONb64) // here is your encoded credentials
	*/

	// tkJSON, err := base64.StdEncoding.DecodeString(cfg.Broker.EncodedCredentials)
	tkJSON, err := base64.StdEncoding.DecodeString(cfg.Broker.EncodedCredentials)
	rawToken := oauth2.Token{}
	if err = json.Unmarshal(tkJSON, &rawToken); err != nil {
		log.Fatal(err)
	}

	cts := config.CachedTokenSource{
		RawToken: rawToken,
	}

	rhClient, rhClientErr := robinhood.Dial(&cts)

	if rhClientErr != nil {
		log.Fatal("Robinhood authentication error %s", rhClientErr)
	}

	commands.InitCommands(rhClient)

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
