package commands

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/weihesdlegend/Mew/utils"

	"github.com/coolboy/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"github.com/weihesdlegend/Mew/clients"
)

// supported commands
var MarketBuyCmd cli.Command
var MarketSellCmd cli.Command
var LimitBuyCmd cli.Command
var LimitSellCmd cli.Command
var AuthCmd cli.Command // auth command

var BufferReader *bufio.Reader

func InitCommands() {
	BufferReader = bufio.NewReader(os.Stdin)
	AuthCmd = cli.Command{
		Name:    "authenticate",
		Aliases: []string{"auth"},
		Usage:   "-u usr -p pwd -m 1111",
		Flags: []cli.Flag{
			&userFlag,
			&passwordFlag,
			&mfaFlag,
		},
		Action: func(ctx *cli.Context) error {
			// TODO move to AuthCommand.go
			log.Info("Creating config file for ", user)
			// Create new config from usr/pwd/mfa
			ts := &robinhood.OAuth{
				Username: user,
				Password: password,
				MFA:      mfa, // Optional
			}

			tk, err := ts.Token()
			if err != nil {
				log.Fatal(err)
			}

			if tk.AccessToken == "" {
				// For some reason the library doesn't return err when password is wrong
				log.Fatal("Auth failed, check your user/password etc...")
			}

			tkJSON, err := json.Marshal(tk)
			tkJSONb64 := base64.StdEncoding.EncodeToString(tkJSON)
			if err != nil {
				log.Fatal(err)
			}

			// .\config.yml
			viper.SetConfigName("config") // name of config file (without extension)
			viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
			viper.AddConfigPath(".")      // optionally look for config in the working directory

			// create empty file if not exist
			file, err := os.OpenFile("./config.yml", os.O_RDONLY|os.O_CREATE, 0666)
			if err != nil {
				log.Fatal(err)
			}
			file.Close() // release it

			viper.Set("broker.name", "robinhood")
			viper.Set("broker.user", user)
			viper.Set("broker.encodedCredentials", tkJSONb64)
			viper.WriteConfig() // Will override

			return nil
		},
	}

	LimitBuyCmd = cli.Command{
		Name:    "limitbuy",
		Aliases: []string{"lb"},
		Usage:   "-t MSFT -l 99.0 -v 2000",
		Flags: []cli.Flag{
			&tickerFlag,
			&sharesFlag,
			&limitBuyFlag,
			&totalValueFlag,
		},
		Action: func(ctx *cli.Context) error {
			rhClient := clients.GetRHClient()

			// init
			lbCmd := &LimitBuyCommand{
				RhClient:     rhClient,
				Ticker:       ticker,
				PercentLimit: limit,
				AmountLimit:  totalValue,
			}

			// Preview
			err := lbCmd.Prepare()
			if err != nil {
				return err
			}

			log.Info(utils.OrderToString(lbCmd.Opts, lbCmd.Ins))

			// Exec
			err = lbCmd.Execute()
			if err != nil {
				log.Error(err)
				return err
			}

			return nil
		},
	}

	// TODO -ls needs to be updated
	LimitSellCmd = cli.Command{
		Name:    "limitsell",
		Aliases: []string{"ls"},
		Usage:   "-t MSFT -ls 101.0 -v 2000",
		Flags: []cli.Flag{
			&tickerFlag,
			&sharesFlag,
			&limitSellFlag,
			&totalValueFlag,
			&percentToSellFlag,
		},
		Action: func(ctx *cli.Context) (err error) {
			rhClient := clients.GetRHClient()

			tickers, tickerParseErr := ParseTicker(ticker)
			if tickerParseErr != nil {
				err = tickerParseErr
				return
			}

			for _, ticker := range tickers {
				// init
				lsCmd := &LimitSellCommand{
					RhClient:     rhClient,
					Ticker:       ticker,
					AmountLimit:  totalValue,
					PercentLimit: limitSell,
				}
				// preview
				if err = lsCmd.Prepare(); err != nil {
					log.Error("Prepare() for ", ticker, " error : ", err)
					continue
				}

				log.Info(utils.OrderToString(lsCmd.Opts, *lsCmd.Ins))

				if err = lsCmd.Execute(); err != nil {
					log.Error("Execute() for ", ticker, " error : ", err)

					continue
				}
			}

			return
		},
	}

	MarketBuyCmd = cli.Command{
		Name:    "buy",
		Aliases: []string{"b"},
		Usage:   "-t MSFT -v 200",
		Flags: []cli.Flag{
			&tickerFlag,
			// TODO &sharesFlag,
			&totalValueFlag,
		},
		Action: func(ctx *cli.Context) error {
			rhClient := clients.GetRHClient()
			// init
			mbCmd := &MarketBuyCommand{
				RhClient:    rhClient,
				Ticker:      ticker,
				AmountLimit: totalValue,
			}

			// Preview
			err := mbCmd.Prepare()
			if err != nil {
				return err
			}
			log.Info(utils.OrderToString(mbCmd.Opts, mbCmd.Ins))

			// Exec
			err = mbCmd.Execute()
			if err != nil {
				log.Error(err)
				return err
			}

			return nil
		},
	}

	MarketSellCmd = cli.Command{
		Name:    "sell",
		Aliases: []string{"s"},
		Usage:   "-t QQQ -v 250",
		Flags: []cli.Flag{
			&tickerFlag,
			// TODO &sharesFlag,
			&totalValueFlag,
			&percentToSellFlag,
		},
		Action: func(ctx *cli.Context) (err error) {
			rhClient := clients.GetRHClient()

			tickers, tickerParseErr := ParseTicker(ticker)
			if tickerParseErr != nil {
				err = tickerParseErr
				return
			}

			for _, ticker := range tickers {
				// init
				msCmd := &MarketSellCommand{
					RhClient:    rhClient,
					Ticker:      ticker,
					AmountLimit: totalValue,
				}
				// preview
				if err = msCmd.Prepare(); err != nil {
					continue
				}

				log.Info(utils.OrderToString(msCmd.Opts, *msCmd.Ins))

				if err = msCmd.Execute(); err != nil {
					continue
				}
			}

			return
		},
	}
}
