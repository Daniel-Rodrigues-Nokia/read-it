// Package version
package version

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"read-it/internal"
	"read-it/internal/components/spinner"
)

const versionCheckCacheTTL = 24 * time.Hour

type githubRelease struct {
	TagName string `json:"tag_name"`
}

type versionCheckCache struct {
	CheckedAt    time.Time `json:"checked_at"`
	RemoteTag    string    `json:"remote_tag"`
	LocalVersion string    `json:"local_version"`
}

type smChecker struct {
	version string
}

var url string = "https://api.github.com/repos/Daniel-Rodrigues-Nokia/read-it/releases/latest"

func NewSmChecker(currentVersion string) smChecker {
	return smChecker{version: currentVersion}
}

func (smC smChecker) CheckVersionWithUI() (bool, error) {
	if needs, ok := smC.tryFromCache(); ok {
		return needs, nil
	}

	spinner, err := spinner.NewSpinner("Checking for updates...", internal.CancelCtrl).Start()
	if err != nil {
		return false, err
	}

	defer func() {
		spinner.Stop()
		internal.ClearStdOut()
	}()

	return smC.checkRemoteVersion()
}

// tryFromCache reports whether the caller should skip the network (second return value).
func (smC smChecker) tryFromCache() (needsUpdated bool, ok bool) {
	path, err := smC.cacheFilePath()
	if err != nil {
		return false, false
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return false, false
	}

	var c versionCheckCache
	if err := json.Unmarshal(data, &c); err != nil {
		return false, false
	}

	if c.LocalVersion != smC.version {
		return false, false
	}
	if time.Since(c.CheckedAt) > versionCheckCacheTTL {
		return false, false
	}

	return smC.version != c.RemoteTag, true
}

func (smC smChecker) cacheFilePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "read-it", "version-check.json"), nil
}

func (smC smChecker) writeCache(remoteTag string) {
	path, err := smC.cacheFilePath()
	if err != nil {
		return
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}

	c := versionCheckCache{
		CheckedAt:    time.Now(),
		RemoteTag:    remoteTag,
		LocalVersion: smC.version,
	}
	data, err := json.Marshal(c)
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0o644)
}

func (smC smChecker) checkRemoteVersion() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tag, err := smC.fetchLatestTag(ctx)
	if err != nil {
		return false, err
	}

	smC.writeCache(tag)

	// TODO: don't compare smv like this
	if smC.version != tag {
		return true, nil
	}

	return false, nil
}

func (smC smChecker) fetchLatestTag(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("version check: unexpected HTTP status from GitHub")
	}

	var content githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
		return "", err
	}

	return content.TagName, nil
}
