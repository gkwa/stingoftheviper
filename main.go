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

var (
	defaultConfigFilename      = "stingoftheviper"
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
	fmt.Fprintln(out, "client.push-frequency", cfg.Client.PushFrequency)
}

func initializeConfig(cmd *cobra.Command, cfg *Config) error {
	v := createViperInstance()

	SetDefaultConfigValues(cfg)

	filePath := "stingoftheviper.yaml"

	// Check if the default config file exists
	if err := WriteDefaultConfigToFile(cfg, filePath); err != nil {
		return err
	}
	fmt.Printf("configfile: %s\n", filePath)

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

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Printf("Error decoding YAML: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("region: %s\n", cfg.SNS.Region)

	// Get the name of the configuration file
	return nil
}
