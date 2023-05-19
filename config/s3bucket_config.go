package config

type S3BucketConfig struct {
	Region string `yaml:"region"`
	Name   string `yaml:"name"`
	S3Path string `yaml:"s3path"`
}
