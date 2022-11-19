package internal

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dhcgn/gh-update/types"
)

var _ WebOperations = (*WebOperationsImpl)(nil)

type WebOperations interface {
	GetGithubRelease(url string) (*[]types.GithubReleaseResult, error)
	GetAssetReader(url string) (data []byte, err error)
}

type WebOperationsImpl struct {
	TestUpdateAssetPath string
}

func (wo WebOperationsImpl) GetAssetReader(url string) (data []byte, err error) {
	if wo.TestUpdateAssetPath != "" {
		return os.ReadFile(wo.TestUpdateAssetPath)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (wo WebOperationsImpl) GetGithubRelease(url string) (*[]types.GithubReleaseResult, error) {
	if wo.TestUpdateAssetPath != "" {
		base := filepath.Base(wo.TestUpdateAssetPath)
		return &[]types.GithubReleaseResult{
			{
				TagName: "v0.0.2",
				Assets: []types.Assets{
					{
						Name:               base,
						BrowserDownloadURL: "https://example.local/" + base,
					},
				},
				PublishedAt: time.Now().AddDate(0, 0, -1),
			},
		}, nil
	}

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github+json")

	if os.Getenv("GITHUB_TOKEN") != "" {
		req.Header.Add("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	ghr := &[]types.GithubReleaseResult{}
	err = json.Unmarshal(body, &ghr)
	if err != nil {
		return nil, err
	}
	return ghr, nil
}
