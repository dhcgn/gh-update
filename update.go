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
	ErrorNoNewVersionFound     = fmt.Errorf("no new version found")
	ErrorLatestNotValid        = fmt.Errorf("latest release not valid")
	ErrorRunningExePathIsEmpty = fmt.Errorf("runningexepath is empty")
)

func SetTestUpdateAssetPath(path string) {
	webop = internal.WebOperationsImpl{
		TestUpdateAssetPath: path,
	}
}

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

type LatestRelease struct {
	Name    string
	Url     string
	Version string
}

// GetLatestVersion get the latest release from github information.
// The latest (newest) release must have a different version than the current version.
// name is the name of the github repository, e.g. "dhcgn/gh-update".
// version is the current version of the application.
// assetfilter is a regex to filter the assets of the release, e.g. "^myapp-.*windows.*zip$".
// The returned LatestRelease contains the name of the asset and the url to download the asset
// and can be used with the func SelfUpdateAndRestart.
func GetLatestVersion(name string, version string, assetfilter string) (LatestRelease, error) {
	// https://api.github.com/repos/dhcgn/workplace-sync/releases
	u, err := url.JoinPath("https://api.github.com/repos/", name, "releases")
	if err != nil {
		return LatestRelease{}, err
	}

	assetRegex, err := regexp.Compile(assetfilter)
	if err != nil {
		return LatestRelease{}, err
	}

	ghr, err := webop.GetGithubRelease(u)
	if err != nil {
		return LatestRelease{}, err
	}

	if len(*ghr) == 0 {
		return LatestRelease{}, fmt.Errorf("no releases found with default branch")
	}

	slices.SortFunc(*ghr, func(i, j types.GithubReleaseResult) bool {
		return i.PublishedAt.After(j.PublishedAt)
	})

	latestRelease := (*ghr)[0]

	if latestRelease.TagName == version {
		return LatestRelease{}, ErrorNoNewVersionFound
	}

	assets := make([]types.Assets, 0)
	for _, asset := range latestRelease.Assets {
		if assetRegex.Match([]byte(asset.Name)) {
			assets = append(assets, asset)
		}
	}

	if len(assets) == 0 {
		return LatestRelease{}, fmt.Errorf("no assets found with filter %s in version %v", assetfilter, latestRelease.TagName)
	}
	if len(assets) > 1 {
		return LatestRelease{}, fmt.Errorf("multiple assets found with filter %s in version %v", assetfilter, latestRelease.TagName)
	}

	return LatestRelease{
		Name:    assets[0].Name,
		Url:     assets[0].BrowserDownloadURL,
		Version: latestRelease.TagName,
	}, nil
}

// SelfUpdateAndRestart updates the current executable with the latest release from github and restarts the application.
// LatestRelease is the result of GetLatestVersion.
// runningexepath is the path to the currently running executable.
func SelfUpdateAndRestart(latest LatestRelease, runningexepath string) error {
	if latest.Version == "" || latest.Url == "" || latest.Name == "" {
		return ErrorLatestNotValid
	}

	if runningexepath == "" {
		return ErrorRunningExePathIsEmpty
	}

	assetData, err := webop.GetAssetReader(latest.Url)
	if err != nil {
		return err
	}

	if strings.HasSuffix(latest.Name, ".zip") {
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

// SelfUpdateWithLatestAndRestart updates the current executable with the latest release from github and restarts the application.
// The latest (newest) release must have a different version than the current version.
// name is the name of the github repository, e.g. "dhcgn/gh-update".
// version is the current version of the application.
// assetfilter is a regex to filter the assets of the release, e.g. "^myapp-.*windows.*zip$".
// runningexepath is the path to the currently running executable.
func SelfUpdateWithLatestAndRestart(name string, version string, assetfilter string, runningexepath string) error {
	latest, err := GetLatestVersion(name, version, assetfilter)
	if err != nil {
		return err
	}

	return SelfUpdateAndRestart(latest, runningexepath)
}
