[![Go Report Card](https://goreportcard.com/badge/github.com/dhcgn/gh-update)](https://goreportcard.com/report/github.com/dhcgn/gh-update)
# gh-update

gh-update is a Go package that provides functionality to update a Go application from a GitHub release.

## Installation

To install gh-update, run the following command:

```bash
go get github.com/dhcgn/gh-update
```

## Usage

### GetLatestVersion

The `GetLatestVersion` function retrieves the latest release from a GitHub repository. It takes the following parameters:

- `name` (string): The name of the GitHub repository, e.g. "dhcgn/gh-update".
- `version` (string): The current version of the application.
- `assetfilter` (string): A regex to filter the assets of the release, e.g. "^myapp-.*windows.*zip$".

The function returns a `LatestRelease` struct, which contains the name of the asset and the URL to download the asset.

### SelfUpdateAndRestart

The `SelfUpdateAndRestart` function updates the current executable with the latest release from GitHub and restarts the application. It takes the following parameters:

- `latest` (LatestRelease): The result of `GetLatestVersion`.
- `runningexepath` (string): The path to the currently running executable.

### SelfUpdateWithLatestAndRestart

The `SelfUpdateWithLatestAndRestart` function updates the current executable with the latest release from GitHub and restarts the application. It takes the following parameters:

- `name` (string): The name of the GitHub repository, e.g. "dhcgn/gh-update".
- `version` (string): The current version of the application.
- `assetfilter` (string): A regex to filter the assets of the release, e.g. "^myapp-.*windows.*zip$".
- `runningexepath` (string): The path to the currently running executable.


## License

This project is licensed under the [MIT License](LICENSE).