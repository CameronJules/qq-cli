/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	verbose     bool
	modelFlag   string
	providerFlg string

	configPaths = []string{
		".",
		"$HOME/.config/qq",
	}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qq [question]", // Note cobra wont parse the brackets so it will run with qq
	Short: "Quick questions for LLMs from your terminal",
	Long:  "qq lets you ask LLMs quick questions from your terminal. Use -v for longer answers. Choose models via -m/--model (e.g., deepseek-chat).",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cobra.MinimumNArgs(1)(cmd, args)
		}
		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: runRootCmd,
}

func runRootCmd(cmd *cobra.Command, args []string) error {
	// Preprocess the question
	question := strings.TrimSpace(strings.Join(args, " "))
	if question == "" {
		return fmt.Errorf("please provide a question, e.g. qq \"What is a goroutine?\"")
	}
	// Select model
	model := strings.TrimSpace(modelFlag)
	if model == "" {
		model = viper.GetString("default_model")
	}
	if model == "" {
		model = "gpt-4o-mini"
	}
	// Select provider
	provider := strings.ToLower(strings.TrimSpace(providerFlg))
	if provider == "" {
		switch {
		case strings.HasPrefix(model, "deepseek-"):
			provider = "deepseek"
		default:
			provider = "openai"
		}
	}
	// Resolve API key (Binding env in viper causes env vars to take precedence)
	var apiKey string
	switch provider {
	case "deepseek":
		apiKey = viper.GetString("deepseek_api_key")
		if apiKey == "" {
			return fmt.Errorf("no DeepSeek API key set. Run `qq configure` or set DEEPSEEK_API_KEY")
		}
	case "openai":
		apiKey = viper.GetString("api_key")
		if apiKey == "" {
			return fmt.Errorf("no OpenAI API key set. Run `qq configure` or set OPENAI_API_KEY")
		}
	default:
		return fmt.Errorf("unknown provider %q (use openai or deepseek)", provider)
	}

	// Initialize client

	// Select the response style

	// Fetch and Print response

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// initialse config
	cobra.OnInitialize(initConfig)

	// Set persistent flags and their defaults
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose (longer) answer")
	rootCmd.PersistentFlags().StringVarP(&modelFlag, "model", "m", "", "Model name (e.g., gpt-4o-mini, deepseek-chat)")
	rootCmd.PersistentFlags().StringVarP(&providerFlg, "provider", "p", "", "Provider (openai or deepseek)")
	// Subcommands
	rootCmd.AddCommand(configureCmd)
	rootCmd.AddCommand(modelsCmd)

}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		for _, p := range configPaths {
			viper.AddConfigPath(os.ExpandEnv(p))
		}
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()
	viper.BindEnv("deepseek_api_key", "DEEPSEEK_API_KEY")
	viper.BindEnv("api_key", "OPENAI_API_KEY")
}
