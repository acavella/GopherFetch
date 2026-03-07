package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

var appVersion = "0000"
var appBuild = "UNK"
var appBuildDate = "00000000"

type Config struct {
	DownloadDir        string `yaml:"download_directory"`
	WorkerLimit        int    `yaml:"worker_limit"`
	RunIntervalSeconds int    `yaml:"run_interval_seconds"`
	Files              []File `yaml:"files"`
}

type File struct {
	URL      string `yaml:"url"`
	Filename string `yaml:"filename"`
}

// loadConfig reads and parses the YAML file
func loadConfig(path string) (Config, error) {
	var config Config
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(yamlFile, &config)
	return config, err
}

func main() {
	// Define the locations to check, in order of priority
	configPaths := []string{
		"gfetch.yaml",
		"/etc/gfetch.yaml",
	}

	var config Config
	var loadedPath string
	var err error

	// 1. Search for and load the config file
	for _, path := range configPaths {
		config, err = loadConfig(path)
		if err == nil {
			loadedPath = path
			fmt.Printf("Loaded configuration from: %s\n", loadedPath)
			break
		}
	}

	if loadedPath == "" {
		log.Fatalf("Failed to load configuration. Searched in: %v", configPaths)
	}

	// Print version and build info at startup
	printver()

	// 2. Run the download process immediately
	executeDownloads(config)

	// 3. Check if we should loop or exit
	if config.RunIntervalSeconds <= 0 {
		fmt.Println("No run_interval_seconds configured (or set to 0). Running once and exiting.")
		return
	}

	// 4. Set up the Ticker for the loop
	ticker := time.NewTicker(time.Duration(config.RunIntervalSeconds) * time.Second)
	defer ticker.Stop()

	fmt.Printf("Started polling every %d seconds. Press Ctrl+C to stop.\n", config.RunIntervalSeconds)

	// 5. The continuous loop
	for range ticker.C {
		// Hot-Reload: Re-read from the path we successfully found earlier
		newConfig, err := loadConfig(loadedPath)
		if err != nil {
			log.Printf("Failed to reload config from %s (skipping this cycle): %v\n", loadedPath, err)
			continue
		}
		
		config = newConfig
		executeDownloads(config)
	}
}

// executeDownloads sets up the worker pool and processes the file queue
func executeDownloads(config Config) {
	if err := os.MkdirAll(config.DownloadDir, os.ModePerm); err != nil {
		log.Printf("Error creating download directory: %v\n", err)
		return
	}

	numWorkers := config.WorkerLimit
	if numWorkers <= 0 {
		numWorkers = 3
	}

	jobs := make(chan File, len(config.Files))
	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, config.DownloadDir, &wg)
	}

	for _, file := range config.Files {
		jobs <- file
	}
	close(jobs)

	wg.Wait()
	fmt.Printf("[%s] Download cycle completed.\n", time.Now().Format("15:04:05"))
	fmt.Println(strings.Repeat("-", 40))
}

// worker constantly pulls from the jobs channel until it is closed and empty
func worker(id int, jobs <-chan File, destDir string, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range jobs {
		destPath := filepath.Join(destDir, file.Filename)
		downloadAndVerifyFile(file.URL, destPath, destDir)
	}
}

// downloadAndVerifyFile handles fetching, hashing, and replacing the file if necessary
func downloadAndVerifyFile(url, destPath, destDir string) {
	tempFile, err := os.CreateTemp(destDir, "temp-dl-*")
	if err != nil {
		log.Printf("Failed to create temp file for %s: %v\n", url, err)
		return
	}
	tempName := tempFile.Name()
	defer os.Remove(tempName)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to download %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Bad status: %s for URL: %s\n", resp.Status, url)
		return
	}

	if _, err = io.Copy(tempFile, resp.Body); err != nil {
		log.Printf("Failed to write data to temp file for %s: %v\n", url, err)
		return
	}
	tempFile.Close()

	newHash, err := hashFile(tempName)
	if err != nil {
		log.Printf("Failed to hash temp file for %s: %v\n", url, err)
		return
	}

	if _, err := os.Stat(destPath); err == nil {
		existingHash, err := hashFile(destPath)
		if err == nil && newHash == existingHash {
			return // Silently skip identical files
		}
	}

	if err := os.Rename(tempName, destPath); err != nil {
		log.Printf("Failed to move temp file to %s: %v\n", destPath, err)
		return
	}

	fmt.Printf("Successfully updated/downloaded: %s\n", destPath)
}

// hashFile generates a SHA-256 hash string for a given file
func hashFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func printver() {
	fmt.Printf("GopherFetch - The Gopher-powered Concurrent File Retrieval Tool\n")
	fmt.Printf("Version: %s\n", appVersion)
	fmt.Printf("Build: %s\n", appBuildDate)
	fmt.Printf("Architecture: %s\n", appBuild)
}