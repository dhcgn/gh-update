package main

import (
	"flag"
	"fmt"
	"os"

	update "github.com/dhcgn/gh-update"
)

var (
	Version = "v0.0.1"
)

var (
	updateFlag     = flag.Bool("update", false, "Check and execute updates")
	updateFileFlag = flag.String("updatefile", "", "Path to update file")
)

func main() {
	fmt.Println("Demo app started", Version)

	if update.IsFirstStartAfterUpdate() {
		fmt.Println("Update finished!")
		oldPid := update.GetOldPid()
		if oldPid != fmt.Sprint(os.Getpid()) {
			err := update.CleanUpAfterUpdate(os.Args[0], oldPid)
			if err != nil {
				fmt.Println("ERROR Clean up:", err)
			}
		} else {
			fmt.Println("ERROR: PID is the same!")
		}
	}

	flag.Parse()

	if *updateFileFlag != "" {
		fmt.Println("Update file:", *updateFileFlag)
		update.SetTestUpdateAssetPath(*updateFileFlag)
	}

	if *updateFlag {
		fmt.Println("Checking for updates ... ")
		err := update.SelfUpdateWithLatestAndRestart("dhcgn/gh-update", Version, "^update_.*exe$", os.Args[0])

		if err != nil && err == update.ErrorNoNewVersionFound {
			fmt.Println("No new version found")
		} else if err != nil {
			fmt.Println("ERROR Update:", err)
		}

		return
	}
	fmt.Println("Exited")
}
