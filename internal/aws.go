package internal

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewAWSConfig(c AwsConfigurations) *aws.Config {
	awsConfig := &aws.Config{
		Region:      aws.String(c.AWS.Region),
		Credentials: credentials.NewStaticCredentials(c.AWS.AccessKey, c.AWS.SecretKey, ""),
	}

	return awsConfig
}

func NewAWSSession(c AwsConfigurations) (*session.Session, error) {
	awsConfig := NewAWSConfig(c)

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}

	return sess, nil
}
