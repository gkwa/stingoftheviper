package config

type SNSConfig struct {
	TopicARN string `yaml:"topic-arn"`
	Region   string `yaml:"region"`
}
