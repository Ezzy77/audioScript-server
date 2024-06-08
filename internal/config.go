package internal

import (
	"errors"
	"github.com/spf13/viper"
)

type AwsConfigurations struct {
	AccessKey    string
	SecretKey    string
	Region       string
	MediaFileUri string
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
	conf.MediaFileUri = viper.GetString("MEDIA_FILE_URI")
	conf.OutputBucket = viper.GetString("OUTPUT_BUCKET")

	// Check that all fields have been populated.
	if conf.AccessKey == "" || conf.SecretKey == "" || conf.Region == "" || conf.MediaFileUri == "" || conf.OutputBucket == "" {
		return nil, errors.New("configuration could not be loaded, check your environment variables")
	}
	return &conf, nil
}
