package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func NewDefaultConfig() *Config {
	return &Config{
		SNS: SNSConfig{
			TopicARN: "arn:aws:sns:us-west-2:123456789012:example-topic",
			Region:   "us-west-2",
		},
		SQS: SQSConfig{
			Region:   "us-west-2",
			QueueARN: "arn:aws:sqs:us-west-2:193048895737",
			QueueURL: "https://sqs.us-west-2.amazonaws.com/193048895737/somename",
		},
		S3Bucket: S3BucketConfig{
			Region: "us-west-2",
			Name:   "mybucket",
			S3Path: ".deliverhalf.yaml",
		},
		Client: ClientConfig{
			PushFrequency: "1m",
		},
	}
}

func ReadConfigFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

var (
	defaultConfigFilename      = "stingoftheviper"
	envPrefix                  = "STING"
	replaceHyphenWithCamelCase = false
)

func main() {
	cmd := NewRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewRootCommand() *cobra.Command {
	var cfg Config

	rootCmd := &cobra.Command{
		Use:   "stingoftheviper",
		Short: "Cober and Viper together at last",
		Long:  `Demonstrate how to get cobra flags to bind to viper properly`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, &cfg)
		},
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()
			PrintConfigValues(out, &cfg)
		},
	}

	return rootCmd
}

func createViperInstance() *viper.Viper {
	v := viper.New()
	v.SetConfigName(defaultConfigFilename)
	v.AddConfigPath(".")
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return v
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := getConfigName(f)
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

func getConfigName(f *pflag.Flag) string {
	configName := f.Name
	if replaceHyphenWithCamelCase {
		configName = strings.ReplaceAll(f.Name, "-", "")
	}
	return configName
}

func PrintConfigValues(out io.Writer, cfg *Config) {
	fmt.Fprintln(out, "Your favorite color is:", cfg.Client.PushFrequency)

	// Print other example fields
	fmt.Fprintln(out, "Example field:", cfg.S3Bucket.Name)
}

// SetDefaultConfigValues sets the default values for the Config struct
func SetDefaultConfigValues(cfg *Config) {
	defaultConfig := NewDefaultConfig()
	cfg.SNS = defaultConfig.SNS
	cfg.SQS = defaultConfig.SQS
	cfg.S3Bucket = defaultConfig.S3Bucket
	cfg.Client = defaultConfig.Client
}

// WriteDefaultConfigToFile writes the default Config struct to a YAML file if the file doesn't exist
func WriteDefaultConfigToFile(cfg *Config, filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		defaultConfigData, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("failed to marshal default config: %v", err)
		}

		err = os.WriteFile(filename, defaultConfigData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write default config file: %v", err)
		}
	}
	return nil
}

func initializeConfig(cmd *cobra.Command, cfg *Config) error {
	v := createViperInstance()

	SetDefaultConfigValues(cfg)
	// Check if the default config file exists
	if err := WriteDefaultConfigToFile(cfg, defaultConfigFilename+".yaml"); err != nil {
		return err
	}

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.AutomaticEnv()
	bindFlags(cmd, v)

	if err := v.Unmarshal(cfg); err != nil {
		return err
	}

	return nil
}
