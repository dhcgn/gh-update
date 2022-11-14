package update

import (
	"testing"
	"time"

	"github.com/dhcgn/gh-update/internal"
	"github.com/dhcgn/gh-update/types"
)

func TestSelfUpdateWithLatestAndRestart(t *testing.T) {

	fops = &FileOperationsMock{}
	osps = &OsOperationsMock{}
	webop = &WebOperationsMock{}

	type args struct {
		name           string
		assetfilter    string
		runningexepath string
		version        string
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		wantErrType error
	}{
		{
			name: "update",
			args: args{
				name:           "owner/repo",
				assetfilter:    "^myapp-.*windows.*zip$",
				runningexepath: `C:\myapp.exe`,
				version:        "v0.0.2",
			},
		},
		{
			name: "no update - same version",
			args: args{
				name:           "owner/repo",
				assetfilter:    "^myapp-.*windows.*zip$",
				runningexepath: `C:\myapp.exe`,
				version:        "v1.2.3",
			},
			wantErr:     true,
			wantErrType: ErrorNoNewVersionFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SelfUpdateWithLatestAndRestart(tt.args.name, tt.args.version, tt.args.assetfilter, tt.args.runningexepath); (err != nil) != tt.wantErr {
				t.Errorf("SelfUpdateWithLatestAndRestart() error = %v, wantErr %v", err, tt.wantErr)

				if tt.wantErrType != nil {
					if err != tt.wantErrType {
						t.Errorf("SelfUpdateWithLatestAndRestart() error = %v, wantErrType %v", err, tt.wantErrType)
					}
				}
			}
		})
	}
}

type FileOperationsMock struct{}

// CleanUpBackup implements internal.FileOperations
func (*FileOperationsMock) CleanUpBackup(p string, try int) error {
	return nil
}

// CreateNewTempPath implements internal.FileOperations
func (*FileOperationsMock) CreateNewTempPath(p string) (newPath string, err error) {
	return internal.FileOperationsImpl{}.CreateNewTempPath(p)
}

// MoveNewExeToOriginalExe implements internal.FileOperations
func (*FileOperationsMock) MoveNewExeToOriginalExe(newPath string, oldPath string) error {
	return nil
}

// MoveRunningExeToBackup implements internal.FileOperations
func (*FileOperationsMock) MoveRunningExeToBackup(p string) error {
	return nil
}

// SaveTo implements internal.FileOperations
func (*FileOperationsMock) SaveTo(data []byte, path string) error {
	return nil
}

// Unzip implements internal.FileOperations
func (*FileOperationsMock) Unzip(zip []byte) (data []byte, err error) {
	return nil, nil
}

type OsOperationsMock struct{}

// Restart implements internal.OsOperations
func (*OsOperationsMock) Restart(path string) error {
	return nil
}

type WebOperationsMock struct{}

// GetAssetReader implements internal.WebOperations
func (*WebOperationsMock) GetAssetReader(url string) (data []byte, err error) {
	return nil, nil
}

// GetGithubRelease implements internal.WebOperations
func (*WebOperationsMock) GetGithubRelease(url string) (*[]types.GithubReleaseResult, error) {
	r := &[]types.GithubReleaseResult{
		{
			TagName:     "v1.2.3",
			PublishedAt: time.Time{},
			Assets: []types.Assets{
				{
					Name:               "myapp-v0.0.3-windows-amd64.zip",
					BrowserDownloadURL: "https://myapp-v0.0.3-windows-amd64.zip",
				},
			},
		},
	}
	return r, nil
}
