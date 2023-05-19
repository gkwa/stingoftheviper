type Config struct {
	SNS      SNSConfig      `yaml:"sns"`
	SQS      SQSConfig      `yaml:"sqs"`
	S3Bucket S3BucketConfig `yaml:"s3bucket"`
	Client   ClientConfig   `yaml:"client"`
}

type ClientConfig struct {
	PushFrequency string `yaml:"push-frequency"`
}

type S3BucketConfig struct {
	Region string `yaml:"region"`
	Name   string `yaml:"name"`
	S3Path string `yaml:"s3path"`
}

type SNSConfig struct {
	TopicARN string `yaml:"topic-arn"`
	Region   string `yaml:"region"`
}

type SQSConfig struct {
	Region   string `yaml:"region"`
	QueueARN string `yaml:"queue-arn"`
	QueueURL string `yaml:"queue-url"`
}
