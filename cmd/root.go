package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "skt",
	Short: "SaaSKit CLI - A helpful CLI tool",
	Long: `SaaSKit CLI is a command line tool that helps you with various tasks.
It provides commands for printing messages and getting help.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`
   _____ _  _______ 
  / ____| |/ /_   _|
 | (___ | ' /  | |  
  \___ \|  <   | |  
  ____) | . \ _| |_ 
 |_____/|_|\_\_____| 
                    
Welcome to SaaSKit CLI!`)
		fmt.Println("Type 'skt init --key <your-key> --' to get started")
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(PrintCmd)
	RootCmd.AddCommand(HelpCmd)
	RootCmd.AddCommand(CloneCmd)
}
