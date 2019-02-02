package main

import (
	"fmt"
	"os"
	"os/exec"
)

// This file is meant to replace the functionality of
// the remoteHandler and the gitHandler. This will rely
// on running command line args to grab the repository
// and check out the correct branch because that will
// improve the performance and prevent copying unnecessary
// packages.

func cloneAndCheckout(build buildInfo, snapshotName string) {
	os.Chdir(getStorageBase())
	os.Mkdir(snapshotName, os.ModeDir)

	uri := ""
	if build.Remote == "" {
		uri = fmt.Sprintf("%s", build.AbsolutePath)
	} else {
		uri = fmt.Sprintf("https://%s@github.com/%s.git", getToken(), build.Remote)
	}

	cmd := exec.Command("git", "clone", uri, snapshotName)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("... When trying to clone %s\n", uri)
	}

	os.Chdir(snapshotName)
	if build.Branch == "" {
		return
	}

	cmd = exec.Command("git", "checkout", build.Branch)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
