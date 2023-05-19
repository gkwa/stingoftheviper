package config

type SQSConfig struct {
	Region   string `yaml:"region"`
	QueueARN string `yaml:"queue-arn"`
	QueueURL string `yaml:"queue-url"`
}
