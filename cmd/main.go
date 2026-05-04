package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"melli_net/internal/config"
	"melli_net/internal/github"
)

var (
	outputDir      string
	outputFilename string
	fileType       string
	noCleanup      bool
	token          string
	repo           string
	maxSize        int
	waitSelector   string
	clickSelector  string
	waitTime       int
	userAgent      string
	timeout        int
)

var rootCmd = &cobra.Command{
	Use:   "melifetch",
	Short: "Download web pages and files via GitHub Actions",
	Long: `melifetch - MeliFetch Network Toolkit
Download web pages and files using GitHub Actions as a proxy.
Supports both simple curl-based fetching and complex browser-based downloads.`,
}

var fetchCmd = &cobra.Command{
	Use:   "fetch [url]",
	Short: "Fetch content using curl (fast, simple)",
	Long: `Fetch content using curl in GitHub Actions.
Best for: direct file downloads, APIs, simple web pages.`,
	Args: cobra.ExactArgs(1),
	RunE: runFetch,
}

var downloadCmd = &cobra.Command{
	Use:   "download [url]",
	Short: "Download using headless browser (complex pages)",
	Long: `Download content using Puppeteer headless browser.
Best for: complex download pages, JavaScript-heavy sites, CurseForge, etc.`,
	Args: cobra.ExactArgs(1),
	RunE: runDownload,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List temporary branches",
	Long:  `List all temporary fetch/download branches in the repository.`,
	RunE:  runList,
}

var cleanCmd = &cobra.Command{
	Use:   "clean [branch-name]",
	Short: "Delete a temporary branch",
	Long:  `Delete a specific temporary branch or all fetch/download branches.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runClean,
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure melifetch settings",
	RunE:  runConfig,
}

func init() {
	// Fetch command flags
	fetchCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")
	fetchCmd.Flags().StringVarP(&outputFilename, "name", "n", "", "Custom output filename")
	fetchCmd.Flags().StringVarP(&fileType, "type", "t", "web", "Content type (web/file)")
	fetchCmd.Flags().IntVar(&maxSize, "max-size", 100, "Max download size in MB")
	fetchCmd.Flags().BoolVar(&noCleanup, "no-cleanup", false, "Keep temporary branch")
	fetchCmd.Flags().StringVar(&token, "token", "", "GitHub token (overrides config)")
	fetchCmd.Flags().StringVar(&repo, "repo", "", "Repository owner/name (overrides config)")
	fetchCmd.Flags().IntVar(&timeout, "timeout", 5, "Timeout in minutes")

	// Download command flags
	downloadCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")
	downloadCmd.Flags().StringVarP(&outputFilename, "name", "n", "", "Custom output filename")
	downloadCmd.Flags().StringVar(&waitSelector, "wait-selector", "", "CSS selector to wait for")
	downloadCmd.Flags().StringVar(&clickSelector, "click-selector", "", "CSS selector to click")
	downloadCmd.Flags().IntVar(&waitTime, "wait-time", 30, "Wait time in seconds")
	downloadCmd.Flags().StringVar(&userAgent, "user-agent", "", "Custom User-Agent")
	downloadCmd.Flags().BoolVar(&noCleanup, "no-cleanup", false, "Keep temporary branch")
	downloadCmd.Flags().StringVar(&token, "token", "", "GitHub token (overrides config)")
	downloadCmd.Flags().StringVar(&repo, "repo", "", "Repository owner/name (overrides config)")
	downloadCmd.Flags().IntVar(&timeout, "timeout", 10, "Timeout in minutes")

	// List command flags
	listCmd.Flags().StringVar(&token, "token", "", "GitHub token (overrides config)")
	listCmd.Flags().StringVar(&repo, "repo", "", "Repository owner/name (overrides config)")

	// Clean command flags
	cleanCmd.Flags().StringVar(&token, "token", "", "GitHub token (overrides config)")
	cleanCmd.Flags().StringVar(&repo, "repo", "", "Repository owner/name (overrides config)")

	// Config command flags
	configCmd.Flags().StringVar(&token, "token", "", "Set GitHub token")
	configCmd.Flags().StringVar(&repo, "repo", "", "Set repository (owner/repo)")

	rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(configCmd)
}

func loadConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if token != "" {
		cfg.Token = token
	}
	if repo != "" {
		cfg.Repo = repo
	}

	return cfg, nil
}

func runConfig(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if token != "" {
		cfg.Token = token
		fmt.Println("✓ Token set")
	}

	if repo != "" {
		cfg.Repo = repo
		fmt.Println("✓ Repo set")
	}

	if token == "" && repo == "" {
		fmt.Println("Current configuration:")
		if cfg.Token != "" {
			fmt.Printf("  Token: %s...%s\n", cfg.Token[:4], cfg.Token[len(cfg.Token)-4:])
		} else {
			fmt.Println("  Token: (not set)")
		}
		if cfg.Repo != "" {
			fmt.Printf("  Repo: %s\n", cfg.Repo)
		} else {
			fmt.Println("  Repo: (not set)")
		}
		path, _ := config.GetConfigPath()
		fmt.Printf("\nConfig file: %s\n", path)
		return nil
	}

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("✓ Configuration saved")
	return nil
}

func runFetch(cmd *cobra.Command, args []string) error {
	url := args[0]

	fmt.Printf("\n%s\n", "═══════════════════════════════════════════════════════════")
	fmt.Printf("🚀 melifetch - Fast Fetch Mode (curl)\n")
	fmt.Printf("%s\n\n", "═══════════════════════════════════════════════════════════")

	fmt.Println("📋 [1/5] Loading configuration...")
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if err := cfg.Validate(); err != nil {
		return err
	}
	fmt.Printf("   ✓ Config loaded: %s\n\n", cfg.Repo)

	fmt.Println("🔗 [2/5] Connecting to GitHub API...")
	fetcher, err := github.NewFetcher(cfg.Token, cfg.Repo)
	if err != nil {
		return fmt.Errorf("❌ Failed to create fetcher: %w", err)
	}
	fmt.Println("   ✓ Connected successfully\n")

	opts := github.FetchOptions{
		URL:            url,
		OutputDir:      outputDir,
		OutputFilename: outputFilename,
		FileType:       fileType,
		MaxSizeMB:      maxSize,
		Cleanup:        !noCleanup,
		Timeout:        time.Duration(timeout) * time.Minute,
	}

	return fetcher.FetchWithOptions(opts)
}

func runDownload(cmd *cobra.Command, args []string) error {
	url := args[0]

	fmt.Printf("\n%s\n", "═══════════════════════════════════════════════════════════")
	fmt.Printf("🚀 melifetch - Browser Download Mode (Puppeteer)\n")
	fmt.Printf("%s\n\n", "═══════════════════════════════════════════════════════════")

	fmt.Println("📋 [1/5] Loading configuration...")
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if err := cfg.Validate(); err != nil {
		return err
	}
	fmt.Printf("   ✓ Config loaded: %s\n\n", cfg.Repo)

	fmt.Println("🔗 [2/5] Connecting to GitHub API...")
	fetcher, err := github.NewFetcher(cfg.Token, cfg.Repo)
	if err != nil {
		return fmt.Errorf("❌ Failed to create fetcher: %w", err)
	}
	fmt.Println("   ✓ Connected successfully\n")

	if waitTime == 0 {
		waitTime = 30
	}

	opts := github.DownloadOptions{
		URL:            url,
		OutputDir:      outputDir,
		OutputFilename: outputFilename,
		WaitSelector:   waitSelector,
		ClickSelector:  clickSelector,
		WaitTime:       waitTime,
		UserAgent:      userAgent,
		Cleanup:        !noCleanup,
		Timeout:        time.Duration(timeout) * time.Minute,
	}

	return fetcher.DownloadWithBrowser(opts)
}

func runList(cmd *cobra.Command, args []string) error {
	fmt.Println("📋 Listing temporary branches...\n")

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	fetcher, err := github.NewFetcher(cfg.Token, cfg.Repo)
	if err != nil {
		return fmt.Errorf("❌ Failed to create fetcher: %w", err)
	}

	fetchBranches, err := fetcher.ListBranches("fetch-")
	if err != nil {
		return fmt.Errorf("❌ Failed to list branches: %w", err)
	}

	downloadBranches, err := fetcher.ListBranches("download-")
	if err != nil {
		return fmt.Errorf("❌ Failed to list branches: %w", err)
	}

	if len(fetchBranches) == 0 && len(downloadBranches) == 0 {
		fmt.Println("No temporary branches found.")
		return nil
	}

	if len(fetchBranches) > 0 {
		fmt.Printf("📥 Fetch branches (%d):\n", len(fetchBranches))
		for _, b := range fetchBranches {
			fmt.Printf("   • %s\n", b.Name)
		}
		fmt.Println()
	}

	if len(downloadBranches) > 0 {
		fmt.Printf("🌐 Download branches (%d):\n", len(downloadBranches))
		for _, b := range downloadBranches {
			fmt.Printf("   • %s\n", b.Name)
		}
		fmt.Println()
	}

	fmt.Printf("Total: %d branches\n", len(fetchBranches)+len(downloadBranches))
	fmt.Println("\nTo delete a branch: melifetch clean <branch-name>")
	fmt.Println("To delete all: melifetch clean")

	return nil
}

func runClean(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	fetcher, err := github.NewFetcher(cfg.Token, cfg.Repo)
	if err != nil {
		return fmt.Errorf("❌ Failed to create fetcher: %w", err)
	}

	if len(args) == 1 {
		branchName := args[0]
		fmt.Printf("🗑️  Deleting branch: %s\n", branchName)
		if err := fetcher.DeleteBranch(branchName); err != nil {
			return err
		}
		fmt.Println("✓ Branch deleted successfully")
		return nil
	}

	fmt.Println("🗑️  Cleaning all temporary branches...\n")

	fetchBranches, _ := fetcher.ListBranches("fetch-")
	downloadBranches, _ := fetcher.ListBranches("download-")

	allBranches := append(fetchBranches, downloadBranches...)

	if len(allBranches) == 0 {
		fmt.Println("No temporary branches to clean.")
		return nil
	}

	deleted := 0
	failed := 0

	for _, b := range allBranches {
		fmt.Printf("Deleting: %s... ", b.Name)
		if err := fetcher.DeleteBranch(b.Name); err != nil {
			fmt.Println("❌ Failed")
			failed++
		} else {
			fmt.Println("✓")
			deleted++
		}
	}

	fmt.Printf("\n✓ Deleted: %d\n", deleted)
	if failed > 0 {
		fmt.Printf("❌ Failed: %d\n", failed)
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "\n❌ Error: %v\n\n", err)
		os.Exit(1)
	}
}
