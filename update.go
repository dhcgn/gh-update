package update

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/dhcgn/gh-update/internal"
	"github.com/dhcgn/gh-update/types"
	"golang.org/x/exp/slices"
)

var (
	fops  internal.FileOperations = internal.FileOperationsImpl{}
	osps  internal.OsOperations   = internal.OsOperationsImpl{}
	webop internal.WebOperations  = internal.WebOperationsImpl{}
)

var (
	ErrorNoNewVersionFound = fmt.Errorf("no new version found")
)

// IsFirstStartAfterUpdate checks if this is the first start after an update
func IsFirstStartAfterUpdate() bool {
	if env := os.Getenv(internal.EnvFinishUpdate); env == "1" {
		return true
	}
	return false
}

// GetOldPid returns the old pid of the application, should be different from the current pid.
func GetOldPid() string {
	return os.Getenv(internal.EnvKillThisPid)
}

// CleanUpAfterUpdate cleans up after an update, should be called after the update.
// It removes the backup of the old executable, executablePath is the path to the new currently running executable.
// A retry is done if the backup file is still in use.
func CleanUpAfterUpdate(executablePath string, oldpid string) error {

	return fops.RemoveExecutable(executablePath, oldpid, 1)
}

// SelfUpdateWithLatestAndRestart updates the current executable with the latest release from github and restarts the application.
// The latest (newest) release must have a different version than the current version.
// name is the name of the github repository, e.g. "dhcgn/gh-update".
// version is the current version of the application.
// assetfilter is a regex to filter the assets of the release, e.g. "^myapp-.*windows.*zip$".
// runningexepath is the path to the currently running executable.
func SelfUpdateWithLatestAndRestart(name string, version string, assetfilter string, runningexepath string) error {
	// https://api.github.com/repos/dhcgn/workplace-sync/releases
	u, err := url.JoinPath("https://api.github.com/repos/", name, "releases")
	if err != nil {
		return err
	}

	assetRegex, err := regexp.Compile(assetfilter)
	if err != nil {
		return err
	}

	if os.Stat(runningexepath); os.IsNotExist(err) {
		return err
	}

	ghr, err := webop.GetGithubRelease(u)
	if err != nil {
		return err
	}

	if len(*ghr) == 0 {
		return fmt.Errorf("no releases found with default branch")
	}

	slices.SortFunc(*ghr, func(i, j types.GithubReleaseResult) bool {
		return i.PublishedAt.After(j.PublishedAt)
	})

	latestRelease := (*ghr)[0]
	// fmt.Println(latestRelease)

	if latestRelease.TagName == version {
		return ErrorNoNewVersionFound
	}

	assets := make([]types.Assets, 0)
	for _, asset := range latestRelease.Assets {
		if assetRegex.Match([]byte(asset.Name)) {
			assets = append(assets, asset)
		}
	}

	if len(assets) == 0 {
		return fmt.Errorf("no assets found with filter %s in version %v", assetfilter, latestRelease.TagName)
	}
	if len(assets) > 1 {
		return fmt.Errorf("multiple assets found with filter %s in version %v", assetfilter, latestRelease.TagName)
	}

	asset := assets[0]
	// fmt.Println("Asset:", asset)

	assetData, err := webop.GetAssetReader(asset.BrowserDownloadURL)
	if err != nil {
		return err
	}

	if strings.HasSuffix(asset.Name, ".zip") {
		assetData, err = fops.Unzip(assetData)
		if err != nil {
			return err
		}
	}

	newpath, err := fops.CreateNewTempPath(runningexepath)
	if err != nil {
		return err
	}

	err = fops.SaveTo(assetData, newpath)
	if err != nil {
		return err
	}
	err = fops.MoveRunningExeToBackup(runningexepath)
	if err != nil {
		return err
	}
	err = fops.MoveNewExeToOriginalExe(newpath, runningexepath)
	if err != nil {
		return err
	}

	err = osps.Restart(runningexepath)
	if err != nil {
		return err
	}

	return nil
}
