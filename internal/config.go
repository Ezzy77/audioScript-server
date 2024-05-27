package internal

import (
	"errors"
	"github.com/spf13/viper"
)

type AwsConfigurations struct {
	AWS struct {
		AccessKey string
		SecretKey string
		Region    string
	}
}

func LoadConfig() (*AwsConfigurations, error) {
	var conf AwsConfigurations

	viper.SetEnvPrefix("aws")
	viper.AutomaticEnv()

	// Set specific variables
	conf.AWS.AccessKey = viper.GetString("ACCESS_KEY_ID")
	conf.AWS.SecretKey = viper.GetString("SECRET_ACCESS_KEY")
	conf.AWS.Region = viper.GetString("REGION")

	// Check that all fields have been populated.
	if conf.AWS.AccessKey == "" || conf.AWS.SecretKey == "" || conf.AWS.Region == "" {
		return nil, errors.New("configuration could not be loaded, check your environment variables")
	}

	return &conf, nil
}
