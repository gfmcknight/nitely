package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Runs a build of the given repository
func buildAction(args argSet) {
	if !args.hasArg(1) {
		fmt.Println("Please specify the item to build.")
		return
	}

	buildInfo := getBuildInfo(nil, args.getArg(1))
	snapshotName := fmt.Sprintf("SNAP-%s-%s",
		buildInfo.Name, time.Now().Format("2006-01-02-150405"))

	ref := ""
	if buildInfo.Branch != "" {
		ref = "refs/heads/" + buildInfo.Branch
	}

	// In the case that we are using a remote, we should update it
	// before trying to build with it
	if buildInfo.Remote != "" {
		cloneOrFetch(*buildInfo)
		ref = "refs/remotes/origin/" + buildInfo.Branch
	}

	inflateRef(ref, buildInfo.AbsolutePath, snapshotName)

	snapshotDir := path.Join(getStorageBase(), snapshotName)
	os.Chdir(snapshotDir)

	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println(err)
		return
	}

	if containsFile(dependencyFile, files) {
		startDependencies(dependencyFile)
		// If we start dependencies, wait some time for them to start
		fmt.Printf("Waiting one minute for services to start...\n")
		time.Sleep(sleepTime(args.getArg(0)))
	}

	for _, filename := range buildFiles {
		if containsFile(filename, files) {
			runner := shell
			if strings.HasSuffix(filename, ".py") {
				runner = "python"
			}

			fmt.Printf("Running file: %s\n", filename)
			cmd := exec.Command(runner, fmt.Sprintf("./%s", filename))

			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println(string(output))
				fmt.Println(err)
			}

			// TODO: Read and compile the test results

			return
		}
	}

	fmt.Printf("No suitable nitely file found to build!\n")

}

// Starts all services in a file, and later will hopefully
// complete a certain set of dependent builds
func startDependencies(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	db := openAndCreateStorage()
	defer db.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		startService(db, scanner.Text())
	}
}

// Starts a specified service
func startService(db *gorm.DB, name string) {
	service := getServiceInfo(db, name)
	cmd := exec.Command(service.AbsolutePath, service.Args)
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
}
