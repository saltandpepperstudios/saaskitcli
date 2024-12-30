package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var CloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a repository",
	Long:  "Clone a repository into a specified directory with authentication",
	Run: func(cmd *cobra.Command, args []string) {
		repoURL, _ := cmd.Flags().GetString("repo")
		path, _ := cmd.Flags().GetString("path")

		// Validate required flags
		if repoURL == "" || path == "" {
			fmt.Println("❌ Required flags are missing:")
			fmt.Println("   --repo: URL of the repository to clone")
			fmt.Println("   --path: Path where to clone the repository")
			return
		}

		// Create absolute path
		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Printf("❌ Error resolving path: %v\n", err)
			return
		}

		// Create directory if it doesn't exist
		err = os.MkdirAll(absPath, 0755)
		if err != nil {
			fmt.Printf("❌ Failed to create directory: %v\n", err)
			return
		}

		fmt.Printf("📁 Created directory: %s\n", absPath)

		// Clone the repository
		fmt.Printf("📥 Cloning repository to %s...\n", absPath)

		gitCmd := exec.Command("git", "clone", repoURL, absPath)
		// Capture output
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr

		if err := gitCmd.Run(); err != nil {
			fmt.Printf("❌ Failed to clone repository: %v\n", err)
			return
		}

		fmt.Println("✅ Repository cloned successfully!")
		fmt.Printf("📝 Next steps:\n")
		fmt.Printf("1. cd %s\n", absPath)
		fmt.Printf("2. Follow the repository's setup instructions\n")
	},
}

func init() {
	CloneCmd.Flags().String("repo", "", "URL of the repository to clone (required)")
	CloneCmd.Flags().String("path", "", "Path where to clone the repository (required)")

	CloneCmd.MarkFlagRequired("repo")
	CloneCmd.MarkFlagRequired("path")
}
