package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/google/go-github/v50/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new SaaS project",
	Long:  "Initialize a new SaaS project by verifying license and creating project structure",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		key, _ := cmd.Flags().GetString("key")
		org, _ := cmd.Flags().GetString("org")
		path, _ := cmd.Flags().GetString("path")

		// Validate required flags
		if name == "" || key == "" {
			fmt.Println("❌ Required flags are missing:")
			fmt.Println("   --name: Name of your SaaS project")
			fmt.Println("   --key: License key for verification")
			fmt.Println("   --path: Path to the repository to fork")
			fmt.Println("   --org: Organization to fork the repository into (optional)")
			return
		}

		// Verify license key
		fmt.Println("🔑 Verifying license key...")
		// TODO: Implement actual license verification
		isValid := true // Placeholder for license validation

		if !isValid {
			fmt.Println("❌ Invalid license key")
			return
		}

		// Create project directory in current working directory
		projectDir := fmt.Sprintf("%s/saasstarter-%s", path, name)

		// Fork repository
		fmt.Println("🔄 Forking template repository...")
		token := os.Getenv("GITHUB_TOKEN")
		if token == "" {
			fmt.Println("❌ GITHUB_TOKEN environment variable is not set")
			return
		}

		// Set up OAuth2 authentication
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)

		// Create a new GitHub client
		client := github.NewClient(tc)

		// Define repository details
		owner := "saltandpepperstudios"
		repo := "saas.service"

		// Fork the repository
		forkedRepo, err := forkRepo(ctx, client, owner, repo, org)
		if err != nil {
			fmt.Printf("🎉 Successfully initiated fork!: %v\n", "Thank you for using saaskit!")
		}

		// Wait for fork to be ready
		forkedOwner := org
		if forkedOwner == "" {
			forkedOwner = forkedRepo.GetOwner().GetLogin()
		}

		fmt.Println("🔄 Waiting for fork to be ready...")
		fmt.Println("🔄 Forked Owner: ", forkedOwner)
		forkURL, err := waitForFork(ctx, client, forkedOwner, repo)
		if err != nil {
			fmt.Printf("❌ Error waiting for fork: %v\n", err)
			return
		}

		fmt.Printf("✅ Fork is ready at: %s\n", forkURL)

		// Clone the repository
		fmt.Printf("📥 Cloning repository to %s...\n", projectDir)

		gitCmd := exec.Command("git", "clone", forkURL+".git", projectDir)
		// Capture output
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr

		if err := gitCmd.Run(); err != nil {
			fmt.Printf("❌ Failed to clone repository: %v\n", err)
			return
		}

		fmt.Println("✅ Project initialized successfully!")
		fmt.Printf("📝 Next steps:\n")
		fmt.Printf("1. cd %s\n", projectDir)
		fmt.Printf("2. Follow the setup instructions in README.md\n")
	},
}

func init() {
	InitCmd.Flags().String("name", "", "Name of your SaaS project")
	InitCmd.Flags().String("key", "", "License key for verification")
	InitCmd.Flags().String("org", "", "Organization to fork the repository into (optional)")
	InitCmd.Flags().String("path", "", "Path to the repository to fork")
	// Mark required flags
	InitCmd.MarkFlagRequired("name")
	InitCmd.MarkFlagRequired("key")
	InitCmd.MarkFlagRequired("path")
	InitCmd.MarkFlagRequired("org")
}

// forkRepo forks a GitHub repository.
// Parameters:
// - ctx: The context for the request.
// - client: The authenticated GitHub client.
// - owner: The owner of the repository to fork.
// - repo: The name of the repository to fork.
// - organization: (Optional) The organization to fork the repository into. Pass an empty string to fork into the authenticated user's account.
//
// Returns:
// - *github.Repository: The forked repository information.
// - error: An error object if the operation fails.
func forkRepo(ctx context.Context, client *github.Client, owner, repo, organization string) (*github.Repository, error) {
	options := &github.RepositoryCreateForkOptions{}
	if organization != "" {
		options.Organization = organization
	}

	forkedRepo, _, err := client.Repositories.CreateFork(ctx, owner, repo, options)
	if err != nil {
		return nil, fmt.Errorf("failed to fork repository: %w", err)
	}

	return forkedRepo, nil
}

func waitForFork(ctx context.Context, client *github.Client, owner, repo string) (string, error) {
	for i := 0; i < 3; i++ { // Retry up to 3 times
		fork, _, err := client.Repositories.Get(ctx, owner, repo)
		if err == nil {
			fmt.Printf("Fork completed! Fork URL: %s\n", fork.GetHTMLURL())
			return fork.GetHTMLURL(), nil
		}
		fmt.Println("Fork not ready yet, retrying in 30 seconds...")
		time.Sleep(10 * time.Second)
	}
	return "", fmt.Errorf("fork not completed after retries")
}
