package config

type Config struct {
	SNS      SNSConfig      `yaml:"sns"`
	SQS      SQSConfig      `yaml:"sqs"`
	S3Bucket S3BucketConfig `yaml:"s3bucket"`
	Client   ClientConfig   `yaml:"client"`
}
