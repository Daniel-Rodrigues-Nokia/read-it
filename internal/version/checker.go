package version

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"read-it/internal"
	"read-it/internal/components/spinner"
)

type GithubRelease struct {
	TagName string `json:"tag_name"`
}

type smChecker struct {
	version string
}

var url string = "https://api.github.com/repos/Daniel-Rodrigues-Nokia/read-it/releases/latest"

func NewSmChecker(currentVersion string) smChecker {
	return smChecker{version: currentVersion}
}

func (smC smChecker) CheckVersionWithUI() (bool, error) {
	spinner, err := spinner.NewSpinner("Checking for updates...", internal.CancelCtrl).Start()
	if err != nil {
		return false, err
	}

	defer func() {
		spinner.Stop()
		internal.ClearStdOut()
	}()

	return smC.checkVersion()
}

func (smC smChecker) checkVersion() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, err
	}

	// do the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// decode data
	var content GithubRelease

	err = json.NewDecoder(resp.Body).Decode(&content)
	if err != nil {
		return false, err
	}

	// compare version
	// TODO: don't compare smv like this
	if smC.version != content.TagName {
		return true, nil
	}

	return false, nil
}
