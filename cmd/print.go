package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var PrintCmd = &cobra.Command{
	Use:   "print",
	Short: "Print a hello message",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from SaaSKit CLI! 👋")
	},
}
