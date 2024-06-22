package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Release struct {
	TagName string `json:"tag_name"`
	Target  string `json:"target_commitish"`
}

// Gets the latest release information from the github repository
func LatestVersion() string {
	url := fmt.Sprint("https://api.github.com/repos/vague2k/rummage/releases/latest")
	resp, err := http.Get(url)
	if err != nil {
		logger.Fatal("Error occured fetching the latest release")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Fatal("Unexpected status code: ", resp.StatusCode)
	}

	var r Release
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		logger.Fatal("Error decoding JSON response: \n", err)
	}

	return fmt.Sprintf("%s from target branch '%s'\n", r.TagName, r.Target)
}
