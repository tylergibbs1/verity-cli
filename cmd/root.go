package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	apiKey  string
	baseURL string
	output  string
)

var rootCmd = &cobra.Command{
	Use:   "verity",
	Short: "Verity CLI - Medicare coverage policies and prior authorization",
	Long: `Verity CLI provides access to Medicare coverage policies, prior authorization
requirements, and medical code lookups from the command line.

Get your API key from: https://verity.backworkai.com/dashboard`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.verity.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Verity API key (or set VERITY_API_KEY env var)")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "https://verity.backworkai.com/api/v1", "API base URL")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "Output format (table, json, yaml)")

	viper.BindPFlag("api_key", rootCmd.PersistentFlags().Lookup("api-key"))
	viper.BindPFlag("base_url", rootCmd.PersistentFlags().Lookup("base-url"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".verity")
	}

	viper.SetEnvPrefix("VERITY")
	viper.AutomaticEnv()

	viper.ReadInConfig()
}

func getAPIKey() string {
	key := viper.GetString("api_key")
	if key == "" {
		fmt.Fprintln(os.Stderr, "Error: API key is required. Set VERITY_API_KEY or use --api-key flag")
		os.Exit(1)
	}
	return key
}

func getBaseURL() string {
	return viper.GetString("base_url")
}

func getOutput() string {
	return viper.GetString("output")
}
