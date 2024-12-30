package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var HelpCmd = &cobra.Command{
	Use:   "help",
	Short: "Get help from the bot",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🤖 Bot Help:")
		fmt.Println("1. Use 'skt print' to print a hello message")
		fmt.Println("2. Use 'skt help' to see this help message")
		fmt.Println("3. Use 'skt init' to initialize a new project")
		fmt.Println("4. Use 'skt clone' to clone a repository")
		fmt.Println("\nFor more detailed information, visit our documentation.")
	},
}
