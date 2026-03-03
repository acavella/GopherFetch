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
	"sync"

	"gopkg.in/yaml.v3"
)

// Config now includes WorkerLimit
type Config struct {
	DownloadDir string `yaml:"download_directory"`
	WorkerLimit int    `yaml:"worker_limit"`
	Files       []File `yaml:"files"`
}

type File struct {
	URL      string `yaml:"url"`
	Filename string `yaml:"filename"`
}

func main() {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config.yaml: %v", err)
	}

	var config Config
	if err = yaml.Unmarshal(yamlFile, &config); err != nil {
		log.Fatalf("Error parsing YAML: %v", err)
	}

	if err = os.MkdirAll(config.DownloadDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating download directory: %v", err)
	}

	// 1. Determine the worker limit (use a fallback if missing or invalid)
	numWorkers := config.WorkerLimit
	if numWorkers <= 0 {
		fmt.Println("Warning: Invalid or missing worker_limit in config. Defaulting to 3.")
		numWorkers = 3
	} else {
		fmt.Printf("Starting worker pool with %d workers.\n", numWorkers)
	}

	// 2. Create a buffered channel to hold our download jobs
	jobs := make(chan File, len(config.Files))
	var wg sync.WaitGroup

	// 3. Start the worker pool using the configured limit
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, config.DownloadDir, &wg)
	}

	// 4. Send all files into the jobs channel
	for _, file := range config.Files {
		jobs <- file
	}
	
	close(jobs)

	// 5. Wait for all workers to finish
	wg.Wait()
	fmt.Println("All downloads processed successfully.")
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
	fmt.Printf("Checking %s...\n", destPath)

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
			fmt.Printf("Skipped: %s (Identical file already exists)\n", destPath)
			return
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