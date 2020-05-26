package commands

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"

	"github.com/coolboy/go-robinhood"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

type AuthCommand struct {
	User     string
	Password string
	MFA      string
}

func (base AuthCommand) Validate() error {
	if base.User == "" {
		return errors.New("User is empty")
	}

	if base.Password == "" {
		return errors.New("Password is empty")
	}

	return nil
}

func (base AuthCommand) Prepare() error {
	return base.Validate()
}

func (base AuthCommand) Execute() error {
	// Create new config from usr/pwd/mfa
	ts := &robinhood.OAuth{
		Username: base.User,
		Password: base.Password,
		MFA:      base.MFA, // Optional
	}

	tk, err := ts.Token()
	if err != nil {
		return err
	}

	if tk.AccessToken == "" {
		// For some reason the library doesn't return err when password is wrong
		return errors.New("auth failed, check your user, password, etc")
	}

	tkJSON, err := json.Marshal(tk)
	tkJSONb64 := base64.StdEncoding.EncodeToString(tkJSON)
	if err != nil {
		return err
	}

	// .\config.yml
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory

	// create empty file if not exist
	file, err := os.OpenFile("./config.yml", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	} // release it

	viper.Set("broker.name", "robinhood")
	viper.Set("broker.user", user)
	viper.Set("broker.encodedCredentials", tkJSONb64)
	if err := viper.WriteConfig(); err != nil {
		return err
	} // Will override

	return nil
}

func AuthCallback(ctx *cli.Context) error {
	log.Info("Creating config file for ", user)

	// init
	authCmd := &AuthCommand{
		User:     user,
		Password: password,
		MFA:      mfa,
	}

	err := authCmd.Prepare()
	if err != nil {
		log.Error(err)
		return err
	}

	err = authCmd.Execute()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil

}
