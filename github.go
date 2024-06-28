package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"fyne.io/fyne/v2"
)

// GetGithubReleasesTags retrieves the tags of all releases from the specified GitHub repository.
// It returns a slice of strings containing the tag names and an error if the request fails.
func GetGithubReleasesTags() ([]string, error) {
	// Specify the URL of the GitHub API endpoint for retrieving repository releases
	url := "https://api.github.com/repos/Taliayaya/Project-86/releases"

	// Send an HTTP GET request to the specified URL
	resp, err := http.Get(url)
	if err != nil {
		// If the request fails, return an error indicating the failure
		return nil, fmt.Errorf("failed to get release information: %v", err)
	}
	defer resp.Body.Close()

	// Check if the HTTP request was successful (status code 200)
	if resp.StatusCode != http.StatusOK {
		// If the request fails, return an error indicating the status code
		return nil, fmt.Errorf("failed to get release information, status code: %d", resp.StatusCode)
	}

	// Decode the response JSON into a slice of structs with a TagName field
	var releases []struct {
		TagName string `json:"tag_name"`
	}

	// Decode the response body into the releases slice
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		// If the decoding fails, return an error indicating the failure
		return nil, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	// Create a slice to store the tag names
	tags := make([]string, len(releases))

	// Iterate over the releases slice and extract the tag names into the tags slice
	for i, release := range releases {
		tags[i] = release.TagName
	}

	// Return the tags slice and nil error indicating success
	fmt.Println("GetGithubReleasesTags success")
	return tags, nil
}

func GetGithubReleasesAssets(tag string) ([]string, map[string]map[string]interface{}, error) {
	// Specify the URL of the GitHub API endpoint for retrieving repository releases
	url := fmt.Sprintf("https://api.github.com/repos/Taliayaya/Project-86/releases/tags/%s", tag)

	// Send an HTTP GET request to the specified URL
	resp, err := http.Get(url)
	if err != nil {
		// If the request fails, return an error indicating the failure
		return nil, nil, fmt.Errorf("failed to get release information: %v", err)
	}
	defer resp.Body.Close()

	// Check if the HTTP request was successful (status code 200)
	if resp.StatusCode != http.StatusOK {
		// If the request fails, return an error indicating the status code
		return nil, nil, fmt.Errorf("failed to get release information, status code: %d", resp.StatusCode)
	}

	// Decode the response JSON into a struct with a Assets field
	var release struct {
		Assets []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
			Size               int64  `json:"size"`
		} `json:"assets"`
	}

	// Decode the response body into the release struct
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		// If the decoding fails, return an error indicating the failure
		return nil, nil, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	// Create a slice to store the asset names
	assets := make([]string, len(release.Assets))

	// Create a map to store the asset names, urls, and sizes
	assetMap := make(map[string]map[string]interface{})

	// Iterate over the release.Assets slice and extract the asset names and urls into the slices and map
	for i, asset := range release.Assets {
		assets[i] = asset.Name
		assetMap[asset.Name] = map[string]interface{}{
			"url":  asset.BrowserDownloadURL,
			"size": asset.Size,
		}
	}

	// Return the asset names slice, asset names url and size map, and nil error indicating success
	fmt.Println("GetGithubReleasesAssets success")
	return assets, assetMap, nil
}

func (l *Launcher) GetGithubReleasesAssetsThread(_app fyne.App, s string) {
	// Runs when the tag changes and fetches asset info
	assets, assetMap, err := GetGithubReleasesAssets(s)
	if err != nil {
		fmt.Println(err)
		_app.SendNotification(fyne.NewNotification("86 Game Launcher: Assets Error", err.Error()))
		_app.Quit()
	}
	l.assets = assets
	l.assetMap = assetMap
	l.assetSelect.SetOptions(l.assets)
	l.assetSelect.ClearSelected()
}
