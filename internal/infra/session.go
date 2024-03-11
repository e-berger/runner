package infra

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

const (
	// AwsRegionKey is the key for the AWS region
	AWSREGIONKEY string = "AWS_REGION"
	// AwsProfileKey is the key for the AWS profile
	AWSPROFILEKEY string = "AWS_PROFILE"
)

var (
	ErrAwsRegionMissing = errors.New("AWS region is missing")
)

func NewSession() (*aws.Config, error) {
	currentRegion := os.Getenv(AWSREGIONKEY)
	if currentRegion == "" {
		return nil, ErrAwsRegionMissing
	}

	return NewSessionForRegion(currentRegion)
}

func NewSessionForRegion(region string) (*aws.Config, error) {
	var err error

	profile := os.Getenv(AWSPROFILEKEY)
	var cfg *aws.Config
	if strings.HasPrefix(region, "http://") {
		cfg, err = createConfigForLocalhost(region, "us-east-1")
	} else {
		if profile != "" {
			cfg, err = createConfigForProfile(region, profile)
		} else {
			cfg, err = createDefaultConfig(region)
		}
	}

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func createDefaultConfig(region string) (*aws.Config, error) {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion(region))
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func createConfigForProfile(region string,
	profile string) (*aws.Config, error) {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(region),
		awsConfig.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func createConfigForLocalhost(awsEndpoint string,
	awsRegion string) (*aws.Config, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if awsEndpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           awsEndpoint,
				SigningRegion: awsRegion,
			}, nil
		}

		// returning EndpointNotFoundError will allow the service to fallback to its default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(awsRegion),
		awsConfig.WithEndpointResolverWithOptions(customResolver),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy")),
	)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
