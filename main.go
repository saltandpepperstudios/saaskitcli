package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "skt",
	Short: "SaaSKit CLI - A helpful CLI tool",
	Long: `SaaSKit CLI is a command line tool that helps you with various tasks.
It provides commands for printing messages and getting help.`,
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print a hello message",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from SaaSKit CLI! 👋")
	},
}

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Get help from the bot",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🤖 Bot Help:")
		fmt.Println("1. Use 'skt print' to print a hello message")
		fmt.Println("2. Use 'skt help' to see this help message")
		fmt.Println("\nFor more detailed information, visit our documentation.")
	},
}

func main() {
	rootCmd.AddCommand(printCmd)
	rootCmd.AddCommand(helpCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
