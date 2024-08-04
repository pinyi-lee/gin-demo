package manager

import (
	"context"

	s3SDK "github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Manager *S3Manager
)

func GetS3() *S3Manager {
	return s3Manager
}

type S3Manager struct {
	client  *s3SDK.Client
	config  S3Config
	context context.Context
}

type S3Config struct {
	AWSS3Region string
	AWSS3Bucket string
}

func (manager *S3Manager) Setup(config S3Config) (err error) {

	return nil
}

func (manager *S3Manager) Close() {

}
