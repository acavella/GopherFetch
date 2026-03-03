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

type Config struct {
	DownloadDir string `yaml:"download_directory"`
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

	var wg sync.WaitGroup

	for _, file := range config.Files {
		wg.Add(1)
		destPath := filepath.Join(config.DownloadDir, file.Filename)
		
		// Pass the DownloadDir as well so we can create temp files on the same drive
		go downloadAndVerifyFile(file.URL, destPath, config.DownloadDir, &wg)
	}

	wg.Wait()
	fmt.Println("All downloads processed successfully.")
}

// downloadAndVerifyFile handles fetching, hashing, and replacing the file if necessary
func downloadAndVerifyFile(url, destPath, destDir string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Checking %s...\n", destPath)

	// 1. Create a temporary file in the destination directory
	tempFile, err := os.CreateTemp(destDir, "temp-dl-*")
	if err != nil {
		log.Printf("Failed to create temp file for %s: %v\n", url, err)
		return
	}
	tempName := tempFile.Name()
	
	// Ensure the temp file is cleaned up if we exit early or after renaming
	defer os.Remove(tempName) 

	// 2. Download the file into the temporary file
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
	
	// We must close the temp file before we can read it to hash it, or rename it (especially on Windows)
	tempFile.Close() 

	// 3. Calculate the hash of the newly downloaded file
	newHash, err := hashFile(tempName)
	if err != nil {
		log.Printf("Failed to hash temp file for %s: %v\n", url, err)
		return
	}

	// 4. Check if the target file exists and compare hashes
	if _, err := os.Stat(destPath); err == nil {
		existingHash, err := hashFile(destPath)
		if err == nil && newHash == existingHash {
			fmt.Printf("Skipped: %s (Identical file already exists)\n", destPath)
			return // The defer os.Remove will silently clean up the temp file
		}
	}

	// 5. Move the temp file to the final destination (overwriting if it exists but differs)
	if err := os.Rename(tempName, destPath); err != nil {
		log.Printf("Failed to move temp file to %s: %v\n", destPath, err)
		return
	}

	fmt.Printf("Successfully updated/downloaded: %s\n", destPath)
}

// hashFile is a helper function that generates a SHA-256 hash string for a given file
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
