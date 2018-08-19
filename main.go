package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

func main() {
	fmt.Println()
	args := os.Args[1:]

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Something went wrong getting the current directory.\n")
		return
	}

	argSet := makeArgSet(args)
	if !argSet.hasArg(0) {
		fmt.Printf("USAGE\n")
		return
	}

	switch argSet.getArg(0) {
	case "add":
		addAction(dir, argSet)
	case "remove":
		removeAction(argSet)
	case "build":
		buildAction(argSet)
	case "list":
		listAction(argSet)
	case "set":
		setAction(argSet)
	}

	fmt.Println()
	fmt.Println(getZen())
}

// Adds a new build or service to the database
func addAction(dir string, args argSet) {
	if args.hasArg(serviceType) {
		addService(dir, args)
		return
	}
	if args.hasArg(remoteType) {
		addFromRemote(args.getArg(nameType), args.getArg(remoteType), args.getArg(branchType))
		return
	}

	if args.hasArg(dirType) {
		newDir, err := filepath.Abs(args.getArg(dirType))
		if err != nil {
			fmt.Println(err)
			return
		}
		dir = newDir
	}

	toAdd := buildInfo{
		AbsolutePath: dir,
		Name:         pathEnd(dir), // Default to the name of the folder
		Branch:       "",           // Use the branch of the working tree
	}

	if args.hasArg(nameType) {
		toAdd.Name = args.getArg(nameType)
	}
	if args.hasArg(branchType) {
		toAdd.Branch = args.getArg(branchType)
	}

	insertBuildInfo(nil, toAdd)
	fmt.Printf("Added build %s on branch %s\nWith path %s\n", toAdd.Name, toAdd.Branch, toAdd.AbsolutePath)
}

// Adds a service that can be started before building
func addService(dir string, args argSet) {
	if args.hasArg(dirType) {
		newDir, err := filepath.Abs(args.getArg(dirType))
		if err != nil {
			fmt.Println(err)
			return
		}
		dir = newDir
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	serviceFilename := args.getArg(serviceType)
	if !containsFile(serviceFilename, files) {
		fmt.Printf("Could not find file %s in directory %s\n", serviceFilename, dir)
		return
	}

	toAdd := serviceInfo{
		Name:         serviceFilename,
		AbsolutePath: path.Join(dir, serviceFilename),
		Args:         "",
	}

	if args.hasArg(nameType) {
		toAdd.Name = args.getArg(nameType)
	}

	if args.hasArg(argsType) {
		toAdd.Args = args.getArg(argsType)
	}

	insertServiceInfo(nil, toAdd)
}

// Removes a build or service from the database
func removeAction(args argSet) {
	if !args.hasArg(1) {
		fmt.Println("Please specify the item to remove from nitely.")
		return
	}

	if args.hasArg(serviceType) {
		deleteServiceInfo(nil, args.getArg(1))
	} else {
		deleteBuildInfo(nil, args.getArg(1))
	}
}

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

// Prints a list of builds to the user that they can run
func listAction(args argSet) {
	if !args.hasArg(1) {
		fmt.Println("Please specify whether to list builds or services.")
		return
	}

	switch args.getArg(1) {

	case "builds":
		fmt.Print("#\tNAME\t\tBRANCH\t\tPATH\t\tREMOTE\n--------------------------------------------\n")
		builds := getBuilds(nil)
		for i := 0; i < len(builds); i++ {
			fmt.Printf("%d.\t%s\t%s\t%s\t%s\n", i+1,
				builds[i].Name, builds[i].Branch, builds[i].AbsolutePath, builds[i].Remote)
		}

	case "services":
		fmt.Print("#\tNAME\t\tARGS\t\tPATH\n--------------------------------------------\n")
		services := getServices(nil)
		for i := 0; i < len(services); i++ {
			fmt.Printf("%d.\t%s\t%s\t%s\n", i+1,
				services[i].Name, services[i].Args, services[i].AbsolutePath)
		}
	}
}

// Set the value of a nitely property
// In the database
func setAction(args argSet) {
	if !args.hasArg(2) {
		fmt.Println("Please the name of the property and the property value.")
		return
	}

	setProperty(nil, args.getArg(1), args.getArg(2))
	fmt.Printf("Property %s set with a value of %s\n", args.getArg(1), args.getArg(2))
}
