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

type SupabaseConfig struct {
	ApiUrl string
	ApiKey string
}

//func LoadSupabaseConfig() (*SupabaseConfig, error) {
//	var conf SupabaseConfig
//	viper.SetEnvPrefix("supabase")
//	viper.AutomaticEnv()
//
//	// Set specific variables
//	conf.ApiUrl = viper.GetString("API_URL")
//	conf.ApiKey = viper.GetString("API_KEY")
//
//	if conf.ApiUrl == "" || conf.ApiKey == "" {
//		return nil, errors.New("Supabase configuration could not be loaded, check your environment variables")
//	}
//	return &conf, nil
//}

func LoadConfig() (*AwsConfigurations, error) {
	var conf AwsConfigurations

	//viper.SetEnvPrefix("aws")
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
