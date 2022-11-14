package update

import (
	"os"
	"strings"
	"testing"
)

func init() {
	if _, err := os.Stat(`C:\dev\githubtoken.txt`); !os.IsNotExist(err) {
		b, err := os.ReadFile(`C:\dev\githubtoken.txt`)
		if err != nil {
			panic(err)
		}
		githubToken := strings.TrimSpace(string(b))
		os.Setenv("GITHUB_TOKEN", githubToken)
	}
}

func TestSelfUpdateWithLatestAndRestart(t *testing.T) {
	type args struct {
		name           string
		assetfilter    string
		runningexepath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestSelfUpdateWithLatestAndRestart",
			args: args{
				name:           "dhcgn/workplace-sync",
				assetfilter:    "^ws-.*zip$",
				runningexepath: `C:\Tools\workplace-sync.exe`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SelfUpdateWithLatestAndRestart(tt.args.name, tt.args.assetfilter, tt.args.runningexepath); (err != nil) != tt.wantErr {
				t.Errorf("SelfUpdateWithLatestAndRestart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
