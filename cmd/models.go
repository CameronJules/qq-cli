package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List common model names and providers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("OpenAI models:")
		fmt.Println("  - gpt-5")

		fmt.Println()
		fmt.Println("DeepSeek models:")
		fmt.Println("  - deepseek-chat        (general chat)")
		fmt.Println("  - deepseek-reasoner    (enhanced reasoning)")
		fmt.Println()
		fmt.Println("Tip: Choose via `-m`, e.g., `qq -m deepseek-chat \"Explain RAG\"`")
	},
}
