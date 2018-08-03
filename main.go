package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	fmt.Println()
	args := os.Args[1:]

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Something went wrong getting the current directory.\n")
		return
	}

	if len(args) < 1 {
		fmt.Printf("USAGE\n")
		return
	}

	switch args[0] {
	case "add":
		addAction(dir, getFlags(args))
	case "remove":
		removeAction(args[1], getFlags(args))
	case "build":
		buildAction(args[1], getFlags(args))
	case "start":
		startAction(nil, args[1])
	case "list":
		listAction(args[1])
	}

	fmt.Println()
	fmt.Println(getZen())
}

// Adds a new build or service to the database
func addAction(dir string, flags map[string]string) {
	if _, exists := flags["-s"]; exists {
		addService(dir, flags)
		return
	}

	if newDir, exists := flags["-d"]; exists {
		newDir, err := filepath.Abs(newDir)
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

	if name, exists := flags["-n"]; exists {
		toAdd.Name = name
	}
	if branch, exists := flags["-b"]; exists {
		toAdd.Branch = branch
	}

	insertBuildInfo(nil, toAdd)
	fmt.Printf("Added build %s on branch %s\nWith path %s\n", toAdd.Name, toAdd.Branch, toAdd.AbsolutePath)
}

// Adds a service that can be started before building
func addService(dir string, flags map[string]string) {
	if newDir, exists := flags["-d"]; exists {
		newDir, err := filepath.Abs(newDir)
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

	serviceFilename := flags["-s"]
	if !containsFile(serviceFilename, files) {
		fmt.Printf("Could not find file %s in directory %s\n", serviceFilename, dir)
		return
	}

	toAdd := serviceInfo{
		Name:         serviceFilename,
		AbsolutePath: path.Join(dir, serviceFilename),
		Args:         "",
	}

	if name, exists := flags["-n"]; exists {
		toAdd.Name = name
	}

	if args, exists := flags["-a"]; exists {
		toAdd.Args = args
	}

	insertServiceInfo(nil, toAdd)
}

// Removes a build or service from the database
func removeAction(name string, flags map[string]string) {
	if _, exists := flags["-s"]; exists {
		deleteServiceInfo(nil, name)
	} else {
		deleteBuildInfo(nil, name)
	}
}

// Runs a build of the given repository
func buildAction(name string, flags map[string]string) {
	buildInfo := getBuildInfo(nil, name)
	snapshotName := fmt.Sprintf("SNAP-%s-%s",
		buildInfo.Name, time.Now().Format("2006-01-02-150405"))

	ref := ""
	if buildInfo.Branch != "" {
		ref = "refs/heads/" + buildInfo.Branch
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
		// This will ideally later be variable/configurable on a project basis
		fmt.Printf("Waiting one minute for services to start...\n")
		time.Sleep(60 * time.Second)
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
		startAction(db, scanner.Text())
	}
}

// Starts a specified service
func startAction(db *sql.DB, name string) {
	service := getServiceInfo(db, name)
	cmd := exec.Command(service.AbsolutePath, service.Args)
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
}

// Prints a list of builds to the user that they can run
func listAction(listType string) {
	switch listType {

	case "builds":
		fmt.Print("#\tNAME\t\tBRANCH\t\tPATH\n--------------------------------------------\n")
		builds := getBuilds(nil)
		for i := 0; i < len(builds); i++ {
			fmt.Printf("%d.\t%s\t%s\t%s\n", i+1,
				builds[i].Name, builds[i].Branch, builds[i].AbsolutePath)
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

// Creates a map of command line flags and their values
// entered by the user.
func getFlags(args []string) map[string]string {
	flags := make(map[string]string)

	for i := 0; i < len(args); i++ {
		if args[i][0] == '-' {
			if i < len(args)-1 {
				flags[args[i]] = args[i+1]
			} else {
				flags[args[i]] = " "
			}
		}
	}

	return flags
}

// Determines whether a directory lists a file under the
// the given name.
func containsFile(name string, files []os.FileInfo) bool {
	for _, file := range files {
		if file.Name() == name {
			return true
		}
	}
	return false
}

// pathEnd gets the final segment of a path.
// For example, /usr/local/bin would yield bin
func pathEnd(dir string) string {
	dir = strings.Replace(dir, "\\", "/", -1)
	tokens := strings.Split(dir, "/")
	return tokens[len(tokens)-1]
}
