/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagOpenAIKey    string
	flagDeepSeekKey  string
	flagDefaultModel string
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Set your API keys and default model",
	Long:  "Configure OpenAI and DeepSeek API keys, and choose a default model used when --model is not provided.",
	RunE:  configure,
}

func configure(cmd *cobra.Command, args []string) error {
	// Load saved keys from Viper
	openAIKey := viper.GetString("open_api_key")
	deepSeekKey := viper.GetString("deepseek_api_key")
	defaultModel := viper.GetString("default_model")

	// Create reader for user input
	reader := bufio.NewReader(os.Stdin)

	// Save new flag values if they have been set, otherwise show key and give user option to type new one
	if flagOpenAIKey == "" {
		fmt.Printf("OpenAI API key [%s]: ", mask(openAIKey))
		in, _ := reader.ReadString('\n')
		in = strings.TrimSpace(in)
		if in != "" {
			openAIKey = in
		}
	} else {
		openAIKey = flagOpenAIKey
	}

	if flagDeepSeekKey == "" {
		fmt.Printf("DeepSeek API key [%s]: ", mask(deepSeekKey))
		in, _ := reader.ReadString('\n')
		in = strings.TrimSpace(in)
		if in != "" {
			deepSeekKey = in
		}
	} else {
		deepSeekKey = flagDeepSeekKey
	}

	// TODO: make this dropdown selectable to avoid invalid values
	if flagDefaultModel == "" {
		fmt.Printf("Default model [%s] (e.g., gpt-5, deepseek-chat): ", defaultModel)
		in, _ := reader.ReadString('\n')
		in = strings.TrimSpace(in)
		if in != "" {
			defaultModel = in
		}
	} else {
		defaultModel = flagDefaultModel
	}

	// Ensure config dir exists
	configDir := os.ExpandEnv("$HOME/.config/qq")
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	cfgPath := filepath.Join(configDir, "config.yaml")

	// Use instance of viper to save values
	v := viper.New()
	if openAIKey != "" {
		v.Set("api_key", openAIKey)
	}
	if deepSeekKey != "" {
		v.Set("deepseek_api_key", deepSeekKey)
	}
	if defaultModel != "" {
		v.Set("default_model", defaultModel)
	}

	// Persist values to config
	v.SetConfigFile(cfgPath)
	v.SetConfigType("yaml")

	if err := v.WriteConfigAs(cfgPath); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Printf("Saved settings to %s\n", cfgPath)
	fmt.Println("Notes:")
	fmt.Println("  - You can override keys with OPENAI_API_KEY and DEEPSEEK_API_KEY.")
	fmt.Println("  - Use `qq -m deepseek-chat \"...\"` to query DeepSeek quickly.")

	return nil
}

func init() {
	// Here you will define your flags and configuration settings.
	configureCmd.Flags().StringVar(&flagOpenAIKey, "openai-key", "", "Set the OpenAI API key")
	configureCmd.Flags().StringVar(&flagDeepSeekKey, "deepseek-key", "", "Set the DeepSeek API key")
	configureCmd.Flags().StringVar(&flagDefaultModel, "default-model", "", "Set the default model (e.g., gpt-4o-mini, deepseek-chat)")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func mask(s string) string {
	if s == "" {
		return ""
	}
	if len(s) <= 6 {
		return "****"
	}
	return s[:3] + "****" + s[len(s)-3:]
}
