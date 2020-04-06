package main

import (
	"os"

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
