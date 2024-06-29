package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"

	"fyne.io/fyne/v2"
)

// DownloadFile downloads a file from GitHub releases and saves it to dest.
func (l *Launcher) DownloadFile(path string, assetDownloadTag string, assetDownloadName string, assetDownloadURL string, assetDownloadSize int64) error {

	fmt.Printf("Downloading asset %s from %s...\n", assetDownloadName, assetDownloadURL)

	// Construct the GitHub Releases API URL
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", "Taliayaya", "Project-86", assetDownloadTag)

	// Make GET request to GitHub Releases API
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get release information: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get release information, status code: %d", resp.StatusCode)
	}

	// Decode the response JSON
	var release struct {
		Assets []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
			Size               int64  `json:"size"`
		} `json:"assets"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("failed to decode JSON response: %v", err)
	}

	// Find the asset by name in the release
	var downloadURL string
	var assetSize int64

	for _, asset := range release.Assets {
		if strings.EqualFold(asset.Name, assetDownloadName) {
			downloadURL = asset.BrowserDownloadURL
			assetSize = asset.Size
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("asset %s not found in release %s", assetDownloadName, l.assetDownloadTag)
	}

	// Check if the file exists
	_, err = os.Stat(path)
	var file *os.File
	var offset int64 = 0 // Offset to resume downloading

	if os.IsNotExist(err) {
		// File does not exist, create a new file
		fmt.Println("File does not exist. Creating a new file.")
		file, err = os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}
	} else {
		// File exists, open it for appending
		fmt.Println("File exists. Resuming download.")
		file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open file for appending: %v", err)
		}

		// Get the size of the existing file
		fi, err := file.Stat()
		if err != nil {
			return fmt.Errorf("failed to get file info: %v", err)
		}
		offset = fi.Size()
		fmt.Printf("File size: %d bytes\n", offset)
	}

	defer file.Close()

	// Create HTTP GET request
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set the Range header to resume download from where it left off
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))

	// Make HTTP request
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Copy data from HTTP response to file with progress indication
	fmt.Println("Resuming download...")

	var written int64 = offset
	progress := make(chan int64)

	go func() {
		var lastReported int64 = -1
		for {
			n, moreErr := io.CopyN(file, resp.Body, 32*1024) // Copy in chunks of 32KB
			if n > 0 {
				written += n
				progress <- written
			}
			if moreErr != nil {
				if moreErr != io.EOF {
					err = moreErr
				}
				break
			}
			if assetSize > 0 {
				currentProgress := written * 100 / assetSize
				if currentProgress != lastReported {
					fmt.Printf("\rProgress: %d%%", currentProgress)

					writtenSize := written
					writtenSizeMB := float64(writtenSize) / (1024 * 1024)
					writtenSizeMBRounded := math.Round(writtenSizeMB*1000) / 1000
					assetSize := assetDownloadSize
					assetSizeMB := float64(assetSize) / (1024 * 1024)
					assetSizeMBRounded := math.Round(assetSizeMB*1000) / 1000
					l.assetDownloadProgress.SetText(fmt.Sprintf("Progress: %0.3f Mb / %0.3f Mb", writtenSizeMBRounded, assetSizeMBRounded))

					lastReported = currentProgress
				}
			}
		}
		close(progress)
	}()

	// Wait for the download to complete
	for p := range progress {
		if assetSize > 0 && p == assetSize {
			fmt.Printf("\rProgress: %d%%\n", p*100/assetSize)

<<<<<<< Updated upstream
			writtenSize := written
			writtenSizeMB := float64(writtenSize) / (1024 * 1024)
			writtenSizeMBRounded := math.Round(writtenSizeMB*1000) / 1000
			assetSize := assetDownloadSize
			assetSizeMB := float64(assetSize) / (1024 * 1024)
			assetSizeMBRounded := math.Round(assetSizeMB*1000) / 1000
			l.assetDownloadProgress.SetText(fmt.Sprintf("Progress: %0.3f Mb / %0.3f Mb", writtenSizeMBRounded, assetSizeMBRounded))
=======
			// writtenSize := written
			// writtenSizeMB := float64(writtenSize) / (1024 * 1024)
			// writtenSizeMBRounded := math.Round(writtenSizeMB*1000) / 1000
			// assetSize := assetDownloadSize
			// assetSizeMB := float64(assetSize) / (1024 * 1024)
			// assetSizeMBRounded := math.Round(assetSizeMB*1000) / 1000
			// l.assetDownloadProgress.SetText(fmt.Sprintf("Progress: %0.3f Mb / %0.3f Mb", writtenSizeMBRounded, assetSizeMBRounded))
>>>>>>> Stashed changes
		}
	}

	if err != nil {
		return fmt.Errorf("error while downloading: %v", err)
	}

	fmt.Printf("\nDownloaded %s to %s\n", assetDownloadName, path)
	return nil
}

func (l *Launcher) DownloadFileThread(_app fyne.App) {
	// Downloads the file to the path specified
	assetTag := l.assetDownloadTag
	assetName := l.assetDownloadName
	assetURL := l.assetDownloadURL
	assetSize := l.assetDownloadSize

	path := fmt.Sprintf("%s/%s", l.assetDownloadDest.Path(), assetName)
	fmt.Println(assetTag, path)

	err := l.DownloadFile(path, assetTag, assetName, assetURL, assetSize)
	if err != nil {
		fmt.Println(err)
		_app.SendNotification(fyne.NewNotification("86 Game Launcher: Download Error", err.Error()))
	}
<<<<<<< Updated upstream
    _app.SendNotification(fyne.NewNotification("86 Game Launcher: Download Status", "Asset has been installed!"))
=======
	_app.SendNotification(fyne.NewNotification("86 Game Launcher: Download Status", "Asset has been installed!"))
>>>>>>> Stashed changes
	l.isDownloading = false
}
