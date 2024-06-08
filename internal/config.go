package internal

import (
	"errors"
	"github.com/spf13/viper"
)

type AwsConfigurations struct {
	AccessKey    string
	SecretKey    string
	Region       string
	UploadBucket string
	OutputBucket string
}

func LoadConfig() (*AwsConfigurations, error) {
	var conf AwsConfigurations

	viper.SetEnvPrefix("aws")
	viper.AutomaticEnv()

	// Set specific variables
	conf.AccessKey = viper.GetString("ACCESS_KEY_ID")
	conf.SecretKey = viper.GetString("SECRET_ACCESS_KEY")
	conf.Region = viper.GetString("REGION")
	conf.UploadBucket = viper.GetString("UPLOAD_BUCKET")
	conf.OutputBucket = viper.GetString("OUTPUT_BUCKET")

	// Check that all fields have been populated.
	if conf.AccessKey == "" || conf.SecretKey == "" || conf.Region == "" || conf.OutputBucket == "" {
		return nil, errors.New("configuration could not be loaded, check your environment variables")
	}
	return &conf, nil
}
