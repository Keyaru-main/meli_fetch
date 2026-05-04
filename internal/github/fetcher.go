package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Fetcher struct {
	Token     string
	RepoOwner string
	RepoName  string
	BaseURL   string
	Client    *http.Client
}

type FetchOptions struct {
	URL            string
	OutputDir      string
	OutputFilename string
	FileType       string
	MaxSizeMB      int
	Cleanup        bool
	Timeout        time.Duration
}

type DownloadOptions struct {
	URL            string
	OutputDir      string
	OutputFilename string
	WaitSelector   string
	ClickSelector  string
	WaitTime       int
	UserAgent      string
	Cleanup        bool
	Timeout        time.Duration
}

type BranchInfo struct {
	Name      string
	CreatedAt time.Time
	Files     []FileInfo
}

type FileInfo struct {
	Name string
	Size int64
	URL  string
}

func NewFetcher(token, repo string) (*Fetcher, error) {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repo format, expected 'owner/repo'")
	}

	return &Fetcher{
		Token:     token,
		RepoOwner: parts[0],
		RepoName:  parts[1],
		BaseURL:   fmt.Sprintf("https://api.github.com/repos/%s/%s", parts[0], parts[1]),
		Client:    &http.Client{Timeout: 180 * time.Second},
	}, nil
}

func (f *Fetcher) doRequestWithRetry(req *http.Request, maxRetries int) (*http.Response, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			wait := time.Duration(i*2) * time.Second
			fmt.Printf("   ⏳ Retry %d/%d (waiting %v)...\n", i+1, maxRetries, wait)
			time.Sleep(wait)
		}

		resp, err := f.Client.Do(req)
		if err == nil {
			if resp.StatusCode == 403 {
				remaining := resp.Header.Get("X-RateLimit-Remaining")
				resetTime := resp.Header.Get("X-RateLimit-Reset")
				if remaining == "0" {
					fmt.Printf("   ⚠️  Rate limit exceeded. Reset at: %s\n", resetTime)
				}
			}
			return resp, nil
		}

		lastErr = err
		fmt.Printf("   ⚠️  Attempt %d failed: %v\n", i+1, err)
	}

	return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}

func (f *Fetcher) TriggerFetchWorkflow(opts FetchOptions) (string, error) {
	fmt.Println("⚙️  [3/5] Triggering GitHub Actions workflow...")

	branchName := fmt.Sprintf("fetch-%d", time.Now().Unix())
	fmt.Printf("   📌 Branch name: %s\n", branchName)
	fmt.Printf("   🌐 Target URL: %s\n", opts.URL)
	fmt.Printf("   📦 Content type: %s\n", opts.FileType)
	if opts.OutputFilename != "" {
		fmt.Printf("   📝 Output filename: %s\n", opts.OutputFilename)
	}
	if opts.MaxSizeMB > 0 {
		fmt.Printf("   📏 Max size: %d MB\n", opts.MaxSizeMB)
	}

	inputs := map[string]string{
		"url":         opts.URL,
		"branch_name": branchName,
		"file_type":   opts.FileType,
	}

	if opts.OutputFilename != "" {
		inputs["output_filename"] = opts.OutputFilename
	}
	if opts.MaxSizeMB > 0 {
		inputs["max_size"] = fmt.Sprintf("%d", opts.MaxSizeMB)
	}

	payload := map[string]interface{}{
		"ref":    "main",
		"inputs": inputs,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	workflowURL := fmt.Sprintf("%s/actions/workflows/fetch.yml/dispatches", f.BaseURL)
	req, err := http.NewRequest("POST", workflowURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "token "+f.Token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	fmt.Println("   🔄 Sending workflow dispatch request...")
	resp, err := f.doRequestWithRetry(req, 3)
	if err != nil {
		return "", fmt.Errorf("❌ Workflow trigger failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("❌ Workflow trigger failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	fmt.Println("   ✓ Workflow triggered successfully\n")
	return branchName, nil
}

func (f *Fetcher) TriggerDownloadWorkflow(opts DownloadOptions) (string, error) {
	fmt.Println("⚙️  [3/5] Triggering download workflow (Puppeteer)...")

	branchName := fmt.Sprintf("download-%d", time.Now().Unix())
	fmt.Printf("   📌 Branch name: %s\n", branchName)
	fmt.Printf("   🌐 Target URL: %s\n", opts.URL)
	if opts.OutputFilename != "" {
		fmt.Printf("   📝 Output filename: %s\n", opts.OutputFilename)
	}
	if opts.WaitSelector != "" {
		fmt.Printf("   🎯 Wait selector: %s\n", opts.WaitSelector)
	}
	if opts.ClickSelector != "" {
		fmt.Printf("   🖱️  Click selector: %s\n", opts.ClickSelector)
	}
	fmt.Printf("   ⏱️  Wait time: %d seconds\n", opts.WaitTime)

	inputs := map[string]string{
		"url":         opts.URL,
		"branch_name": branchName,
		"wait_time":   fmt.Sprintf("%d", opts.WaitTime),
	}

	if opts.OutputFilename != "" {
		inputs["output_filename"] = opts.OutputFilename
	}
	if opts.WaitSelector != "" {
		inputs["wait_selector"] = opts.WaitSelector
	}
	if opts.ClickSelector != "" {
		inputs["click_selector"] = opts.ClickSelector
	}
	if opts.UserAgent != "" {
		inputs["user_agent"] = opts.UserAgent
	}

	payload := map[string]interface{}{
		"ref":    "main",
		"inputs": inputs,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	workflowURL := fmt.Sprintf("%s/actions/workflows/download.yml/dispatches", f.BaseURL)
	req, err := http.NewRequest("POST", workflowURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "token "+f.Token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	fmt.Println("   🔄 Sending workflow dispatch request...")
	resp, err := f.doRequestWithRetry(req, 3)
	if err != nil {
		return "", fmt.Errorf("❌ Workflow trigger failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("❌ Workflow trigger failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	fmt.Println("   ✓ Workflow triggered successfully\n")
	return branchName, nil
}

func (f *Fetcher) WaitForBranch(branchName string, timeout time.Duration) error {
	fmt.Println("⏳ [4/5] Waiting for workflow to complete...")
	fmt.Printf("   ⏱️  Timeout: %v\n", timeout)
	fmt.Printf("   🔍 Checking branch: %s\n", branchName)

	start := time.Now()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	attempt := 0
	lastStatus := ""

	for {
		select {
		case <-ticker.C:
			attempt++
			elapsed := time.Since(start).Round(time.Second)

			branchURL := fmt.Sprintf("%s/branches/%s", f.BaseURL, branchName)
			req, _ := http.NewRequest("GET", branchURL, nil)
			req.Header.Set("Authorization", "token "+f.Token)
			req.Header.Set("Accept", "application/vnd.github.v3+json")

			resp, err := f.Client.Do(req)
			if err == nil && resp.StatusCode == 200 {
				resp.Body.Close()
				fmt.Printf("\n   ✓ Branch created after %v (%d checks)\n\n", elapsed, attempt)
				return nil
			}
			if resp != nil {
				resp.Body.Close()
			}

			spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
			status := fmt.Sprintf("   %s Waiting... [%v elapsed, attempt %d]", spinner[attempt%len(spinner)], elapsed, attempt)
			if status != lastStatus {
				fmt.Printf("\r%s", status)
				lastStatus = status
			}

		case <-time.After(timeout):
			fmt.Println()
			return fmt.Errorf("❌ Timeout after %v", timeout)
		}

		if time.Since(start) > timeout {
			fmt.Println()
			return fmt.Errorf("❌ Timeout after %v", timeout)
		}
	}
}

func (f *Fetcher) DownloadContent(branchName, outputDir string) error {
	fmt.Println("📥 [5/5] Downloading content...")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("❌ Failed to create output directory: %w", err)
	}
	fmt.Printf("   📁 Output directory: %s\n", outputDir)

	fmt.Println("\n   📋 Fetching metadata...")
	metadataURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/fetched/metadata.txt",
		f.RepoOwner, f.RepoName, branchName)

	resp, err := http.Get(metadataURL)
	if err == nil && resp.StatusCode == 200 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Printf("\n%s\n", strings.Repeat("─", 60))
		fmt.Printf("%s", string(body))
		fmt.Printf("%s\n\n", strings.Repeat("─", 60))
	} else {
		fmt.Println("   ⚠️  No metadata found")
	}

	fmt.Println("   🔍 Getting filename...")
	filenameURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/fetched/filename.txt",
		f.RepoOwner, f.RepoName, branchName)

	resp, err = http.Get(filenameURL)
	var filename string
	if err == nil && resp.StatusCode == 200 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		filename = strings.TrimSpace(string(body))
		fmt.Printf("   📄 Filename: %s\n", filename)
	}

	if filename == "" {
		fmt.Println("   ⚠️  No filename found, trying defaults...")
		files := []string{
			"content.html", "content.json", "content.bin", "content.txt",
			"content.pdf", "content.zip", "content.jpg", "content.png",
			"content.jar", "content.exe", "content.apk", "content.tar.gz",
			"content.mp4", "content.mp3", "content.xml", "content.csv",
			"content.7z", "content.rar", "content.gz", "content.tar",
			"content.mkv", "content.webp", "content.gif",
			"content",
		}

		for _, file := range files {
			testURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/fetched/%s",
				f.RepoOwner, f.RepoName, branchName, file)
			resp, err := http.Head(testURL)
			if err == nil && resp.StatusCode == 200 {
				resp.Body.Close()
				filename = file
				fmt.Printf("   ✓ Found: %s\n", filename)
				break
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
	}

	if filename == "" {
		return fmt.Errorf("❌ Could not determine filename. Branch may not have completed successfully")
	}

	fmt.Printf("\n   📥 Downloading: %s\n", filename)
	fileURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/fetched/%s",
		f.RepoOwner, f.RepoName, branchName, filename)

	startTime := time.Now()
	resp, err = http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("❌ Failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("❌ Failed to download file (HTTP %d)", resp.StatusCode)
	}

	outputPath := filepath.Join(outputDir, filename)
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("❌ Failed to create output file: %w", err)
	}
	defer outFile.Close()

	written, err := f.copyWithProgress(outFile, resp.Body, resp.ContentLength)
	if err != nil {
		return fmt.Errorf("❌ Failed to write file: %w", err)
	}

	elapsed := time.Since(startTime)
	speed := float64(written) / elapsed.Seconds() / 1024 / 1024

	fmt.Printf("\n   ✓ Downloaded: %s\n", outputPath)
	fmt.Printf("   📊 Size: %s\n", formatBytes(written))
	fmt.Printf("   ⚡ Speed: %.2f MB/s\n", speed)
	fmt.Printf("   ⏱️  Time: %v\n", elapsed.Round(time.Millisecond))

	f.downloadOptionalFiles(branchName, outputDir)

	return nil
}

func (f *Fetcher) copyWithProgress(dst io.Writer, src io.Reader, totalSize int64) (int64, error) {
	var written int64
	buf := make([]byte, 32*1024)
	lastPrint := time.Now()

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				return written, ew
			}
			if nr != nw {
				return written, io.ErrShortWrite
			}

			if time.Since(lastPrint) > 500*time.Millisecond {
				if totalSize > 0 {
					percent := float64(written) / float64(totalSize) * 100
					fmt.Printf("\r   📥 Progress: %.1f%% (%s / %s)", percent, formatBytes(written), formatBytes(totalSize))
				} else {
					fmt.Printf("\r   📥 Downloaded: %s", formatBytes(written))
				}
				lastPrint = time.Now()
			}
		}
		if er != nil {
			if er != io.EOF {
				return written, er
			}
			break
		}
	}
	fmt.Println()
	return written, nil
}

func (f *Fetcher) downloadOptionalFiles(branchName, outputDir string) {
	optionalFiles := []string{
		"metadata.json",
		"page_initial.png",
		"page_final.png",
		"found_links.json",
	}

	for _, file := range optionalFiles {
		fileURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/fetched/%s",
			f.RepoOwner, f.RepoName, branchName, file)

		resp, err := http.Get(fileURL)
		if err == nil && resp.StatusCode == 200 {
			outputPath := filepath.Join(outputDir, file)
			outFile, err := os.Create(outputPath)
			if err == nil {
				io.Copy(outFile, resp.Body)
				outFile.Close()
				fmt.Printf("   📎 Downloaded: %s\n", file)
			}
		}
		if resp != nil {
			resp.Body.Close()
		}
	}
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func (f *Fetcher) DeleteBranch(branchName string) error {
	fmt.Println("\n🗑️  Cleaning up temporary branch...")

	deleteURL := fmt.Sprintf("%s/git/refs/heads/%s", f.BaseURL, branchName)
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+f.Token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := f.Client.Do(req)
	if err != nil {
		return fmt.Errorf("⚠️  Could not delete branch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		fmt.Printf("   ✓ Branch '%s' deleted\n", branchName)
		return nil
	}

	return fmt.Errorf("⚠️  Could not delete branch (HTTP %d)", resp.StatusCode)
}

func (f *Fetcher) ListBranches(prefix string) ([]BranchInfo, error) {
	fmt.Printf("🔍 Listing branches with prefix: %s\n", prefix)

	branchesURL := fmt.Sprintf("%s/branches", f.BaseURL)
	req, err := http.NewRequest("GET", branchesURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+f.Token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := f.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to list branches (HTTP %d)", resp.StatusCode)
	}

	var branches []struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&branches); err != nil {
		return nil, err
	}

	var result []BranchInfo
	for _, b := range branches {
		if strings.HasPrefix(b.Name, prefix) {
			result = append(result, BranchInfo{Name: b.Name})
		}
	}

	return result, nil
}

func (f *Fetcher) Fetch(url, outputDir, fileType string, cleanup bool) error {
	opts := FetchOptions{
		URL:       url,
		OutputDir: outputDir,
		FileType:  fileType,
		Cleanup:   cleanup,
		Timeout:   5 * time.Minute,
	}
	return f.FetchWithOptions(opts)
}

func (f *Fetcher) FetchWithOptions(opts FetchOptions) error {
	startTime := time.Now()

	branchName, err := f.TriggerFetchWorkflow(opts)
	if err != nil {
		return err
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 5 * time.Minute
	}

	if err := f.WaitForBranch(branchName, timeout); err != nil {
		return err
	}

	time.Sleep(3 * time.Second)

	if err := f.DownloadContent(branchName, opts.OutputDir); err != nil {
		if opts.Cleanup {
			f.DeleteBranch(branchName)
		}
		return err
	}

	if opts.Cleanup {
		f.DeleteBranch(branchName)
	}

	elapsed := time.Since(startTime).Round(time.Second)
	fmt.Printf("\n%s\n", "═══════════════════════════════════════════════════════════")
	fmt.Printf("✅ Fetch completed successfully in %v\n", elapsed)
	fmt.Printf("%s\n\n", "═══════════════════════════════════════════════════════════")

	return nil
}

func (f *Fetcher) DownloadWithBrowser(opts DownloadOptions) error {
	startTime := time.Now()

	branchName, err := f.TriggerDownloadWorkflow(opts)
	if err != nil {
		return err
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 10 * time.Minute
	}

	if err := f.WaitForBranch(branchName, timeout); err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	if err := f.DownloadContent(branchName, opts.OutputDir); err != nil {
		if opts.Cleanup {
			f.DeleteBranch(branchName)
		}
		return err
	}

	if opts.Cleanup {
		f.DeleteBranch(branchName)
	}

	elapsed := time.Since(startTime).Round(time.Second)
	fmt.Printf("\n%s\n", "═══════════════════════════════════════════════════════════")
	fmt.Printf("✅ Download completed successfully in %v\n", elapsed)
	fmt.Printf("%s\n\n", "═══════════════════════════════════════════════════════════")

	return nil
}
