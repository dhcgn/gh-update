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

func IsUpdateFinished() bool {
	if env := os.Getenv(internal.EnvFinishUpdate); env == "1" {
		return true
	}
	return false
}

func CleanUp(executablePath string) error {
	return fops.CleanUpBackup(executablePath, 1)
}

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
		return fmt.Errorf("no assets found with filter %s", assetfilter)
	}
	if len(assets) > 1 {
		return fmt.Errorf("multiple assets found with filter %s", assetfilter)
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
