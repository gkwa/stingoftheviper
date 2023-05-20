package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	envPrefix                  = "STING"
	replaceHyphenWithCamelCase = false
)

func main() {
	cmd := NewRootCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

func NewRootCommand() *cobra.Command {
	var cfg Config

	rootCmd := &cobra.Command{
		Use:   "stingoftheviper",
		Short: "Cobra and Viper together at last",
		Long:  "Demonstrate how to get cobra flags to bind to viper properly",
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

func createViperInstance(configFilename string) *viper.Viper {
	v := viper.New()
	v.SetConfigName(configFilename)
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
	fmt.Fprintln(out, "s3bucket.name", cfg.S3Bucket.Name)
	fmt.Fprintln(out, "client.pushfrequency", cfg.Client.PushFrequency)
	fmt.Fprintln(out, "sqs.region", cfg.SNS.Region)
}

func initializeConfig(cmd *cobra.Command, cfg *Config) error {
	defaultConfigFilename := "stingoftheviper"
	v := createViperInstance(defaultConfigFilename)

	SetDefaultConfigValues(cfg)

	// Check if the default config file exists
	if err := WriteDefaultConfigToFile(cfg, defaultConfigFilename+".yaml"); err != nil {
		return err
	}

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		fmt.Print(err)
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
