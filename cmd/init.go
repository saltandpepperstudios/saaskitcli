package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/google/go-github/v50/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// LicenseResponse represents the response from the license validation API
type LicenseResponse []struct {
	Key string `json:"key"`
	// Add other fields if needed
}

// ValidateKey checks if the license key is valid by calling the SheetDB API
func ValidateKey(key string) (bool, error) {
	// Create the API URL with the key
	apiURL := fmt.Sprintf("https://sheetdb.io/api/v1/ju2p5lmgeed0j/search?key=%s", key)

	// Create a new request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return false, fmt.Errorf("error creating request: %v", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", "Bearer 7nn8qwhuvlmr4t34xhc16jla3g73xzsz3p463k37")

	// Make the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response: %v", err)
	}

	// Parse the response
	var licenseResp LicenseResponse
	if err := json.Unmarshal(body, &licenseResp); err != nil {
		return false, fmt.Errorf("error parsing response: %v", err)
	}

	// Check if the response is empty
	if len(licenseResp) == 0 {
		return false, fmt.Errorf("invalid license key. Please contact support@saasstarter.live")
	}

	return true, nil
}

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
		isValid, err := ValidateKey(key)
		if err != nil {
			fmt.Printf("❌ License validation failed: %v\n", err)
			return
		}

		if !isValid {
			fmt.Println("❌ Invalid license key")
			fmt.Println("Please contact support@saasstarter.live to obtain a valid license key")
			return
		}

		fmt.Println("✅ License key verified successfully!")

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
		frontendRepo := "starterkit.client"

		// Fork the backend repository
		forkedRepo, err := forkRepo(ctx, client, owner, repo, org)
		if err != nil {
			fmt.Printf("🎉 Successfully initiated fork for Backend!: %v\n", "Thank you for using GO BACKEND for SaaS Starter Kit!")
		}

		// Fork the frontend repository
		_, err = forkRepo(ctx, client, owner, frontendRepo, org)
		if err != nil {
			fmt.Printf("🎉 Successfully initiated fork for Frontend!: %v\n", "Thank you for using GO FRONTEND for SaaS Starter Kit!")
		}

		// Wait for fork to be ready
		forkedOwner := org
		if forkedOwner == "" {
			forkedOwner = forkedRepo.GetOwner().GetLogin()
		}

		fmt.Println("🔄 Waiting for fork to be ready...")
		fmt.Println("🔄 Forked Owner: ", forkedOwner)
		forkURL, frontendForkURL, err := waitForFork(ctx, client, forkedOwner, repo, frontendRepo)
		if err != nil {
			fmt.Printf("❌ Error waiting for fork: %v\n", err)
			return
		}

		fmt.Printf("✅ Fork is ready at: %s\n", forkURL)

		// Clone the repository
		fmt.Printf("📥 Cloning repository to %s...\n", projectDir)

		// Clone backend repository
		backendDir := projectDir + "/backend"
		gitCmd := exec.Command("git", "clone", forkURL+".git", backendDir)
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr
		if err := gitCmd.Run(); err != nil {
			fmt.Printf("❌ Failed to clone backend repository: %v\n", err)
			return
		}

		// Clone frontend repository
		frontendDir := projectDir + "/frontend"
		frontendGitCmd := exec.Command("git", "clone", frontendForkURL+".git", frontendDir)
		frontendGitCmd.Stdout = os.Stdout
		frontendGitCmd.Stderr = os.Stderr
		if err := frontendGitCmd.Run(); err != nil {
			fmt.Printf("❌ Failed to clone frontend repository: %v\n", err)
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

func waitForFork(ctx context.Context, client *github.Client, owner, repo, frontendRepo string) (string, string, error) {
	for i := 0; i < 3; i++ { // Retry up to 3 times
		fork, _, _ := client.Repositories.Get(ctx, owner, repo)
		frontendFork, _, err := client.Repositories.Get(ctx, owner, frontendRepo)
		if err == nil {
			fmt.Printf("Fork completed! Fork URL: %s\n", fork.GetHTMLURL())
			return fork.GetHTMLURL(), frontendFork.GetHTMLURL(), nil
		}
		fmt.Println("Fork not ready yet, retrying in 30 seconds...")
		time.Sleep(10 * time.Second)
	}
	return "", "", fmt.Errorf("fork not completed after retries")
}
