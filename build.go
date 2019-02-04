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

	db := openAndCreateStorage()

	buildInfo := getBuildInfo(db, args.getArg(1))
	snapshotName := fmt.Sprintf("SNAP-%s-%s",
		buildInfo.Name, time.Now().Format("2006-01-02-150405"))

	cloneAndCheckout(*buildInfo, snapshotName)

	snapshotDir := path.Join(getStorageBase(), snapshotName)
	os.Chdir(snapshotDir)

	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println(err)
		return
	}

	if containsFile(dependencyFile, files) {
		fmt.Println("Starting services")
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
			cmd := exec.Command(runner, fmt.Sprintf("./%s %s", filename, buildInfo.AbsolutePath))

			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println(string(output))
				fmt.Println(err)
			}

			if _, err := os.Stat(resultsFile); os.IsNotExist(err) {
				fmt.Printf("No file %s to take results from!", resultsFile)
				return
			}

			saveResults(db, buildInfo, resultsFile)

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

func saveResults(db *gorm.DB, buildInfo *buildInfo, filename string) {
	if db == nil {
		db = openAndCreateStorage()
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	run := testRun{
		DateRun: time.Now(),
		BuildID: buildInfo.ID,
		Results: make([]testResult, 0),
	}

	passedAll := true
	for scanner.Scan() {
		run.Results = append(run.Results, testResult{
			Passed:   scanner.Text()[0] == 'P',
			TestName: scanner.Text()[1:],
		})

		if !run.Results[len(run.Results)-1].Passed {
			passedAll = false
		}
	}

	if passedAll {
		setStatus(buildInfo.Remote, "success", getCommitID())
	} else {
		setStatus(buildInfo.Remote, "failure", getCommitID())
	}

	insertTestRun(db, run)
}
