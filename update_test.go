package update

import (
	"reflect"
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

// RemoveExecutable implements internal.FileOperations
func (*FileOperationsMock) RemoveExecutable(p string, pid string, try int) error {
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

func TestGetLatestVersion(t *testing.T) {

	fops = &FileOperationsMock{}
	osps = &OsOperationsMock{}
	webop = &WebOperationsMock{}

	type args struct {
		name        string
		version     string
		assetfilter string
	}
	tests := []struct {
		name    string
		args    args
		want    LatestRelease
		wantErr bool
	}{
		{
			name: "get latest version",
			args: args{
				name:        "owner/repo",
				assetfilter: "^myapp-.*windows.*zip$",
				version:     "v0.0.2",
			},
			want: LatestRelease{
				Name:    "myapp-v0.0.3-windows-amd64.zip",
				Url:     `https://myapp-v0.0.3-windows-amd64.zip`,
				Version: "v1.2.3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLatestVersion(tt.args.name, tt.args.version, tt.args.assetfilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLatestVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelfUpdateAndRestart(t *testing.T) {

	fops = &FileOperationsMock{}
	osps = &OsOperationsMock{}
	webop = &WebOperationsMock{}

	type args struct {
		latest         LatestRelease
		runningexepath string
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		wantErrorType error
	}{
		{
			name: "Empty LatestRelease",
			args: args{
				latest:         LatestRelease{},
				runningexepath: "",
			},
			wantErr:       true,
			wantErrorType: ErrorLatestNotValid,
		},
		{
			name: "Set LatestRelease but ErrorRunningExePathIsEmpty",
			args: args{
				latest: LatestRelease{
					Name:    "myapp-v0.0.3-windows-amd64.zip",
					Url:     `https://myapp-v0.0.3-windows-amd64.zip`,
					Version: "v1.2.3",
				},
				runningexepath: "",
			},
			wantErr:       true,
			wantErrorType: ErrorRunningExePathIsEmpty,
		},
		{
			name: "Not Error",
			args: args{
				latest: LatestRelease{
					Name:    "myapp-v0.0.3-windows-amd64.zip",
					Url:     `https://myapp-v0.0.3-windows-amd64.zip`,
					Version: "v1.2.3",
				},
				runningexepath: "myapp.exe",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SelfUpdateAndRestart(tt.args.latest, tt.args.runningexepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelfUpdateAndRestart() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && tt.wantErrorType != nil && err != tt.wantErrorType {
				t.Errorf("SelfUpdateAndRestart() error = %v, wantErrorType %v", err, tt.wantErrorType)
			}
		})
	}
}
