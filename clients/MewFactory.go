package clients

import (
	"encoding/base64"
	"encoding/json"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/weihesdlegend/Mew/config"
	"golang.org/x/oauth2"
)

// Global singleton instances
var gRHClient *RHClient
var gAppConfig config.Configurations

// once protectors
var cfgOnce sync.Once
var cliOnce sync.Once

func GetRHClient() *RHClient {
	cliOnce.Do(func() {
		log.Info("Creating rhClient...")

		cfg := GetAppConfig()

		tkJSON, err := base64.StdEncoding.DecodeString(cfg.Broker.EncodedCredentials)
		rawToken := oauth2.Token{}
		if err = json.Unmarshal(tkJSON, &rawToken); err != nil {
			log.Fatal(err)
		}

		cts := config.CachedTokenSource{
			RawToken: rawToken,
		}

		gRHClient = &RHClient{}
		rhClientErr := gRHClient.Init(&cts)
		if rhClientErr != nil {
			log.Fatal(rhClientErr)
		}

	})

	return gRHClient
}

func GetAppConfig() *config.Configurations {
	cfgOnce.Do(func() {
		log.Info("Loading config...")
		viper.SetConfigName("config") // name of config file (without extension)
		viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath(".")      // optionally look for config in the working directory

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Config read error! %v", err)
		}

		if err := viper.Unmarshal(&gAppConfig); err != nil {
			log.Fatalf("Unable to decode into struct, %v", err)
		}
	})

	return &gAppConfig
}
