package main

import (
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
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

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
	case "list":
		listAction()
	}
}

// Adds a new build or service to the database
func addAction(dir string, flags map[string]string) {
	if newDir, exists := flags["-d"]; exists {
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

func removeAction(name string, flags map[string]string) {

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

	inflateCommit(ref, buildInfo.AbsolutePath, snapshotName)

	snapshotDir := path.Join(getStorageBase(), snapshotName)
	os.Chdir(snapshotDir)

	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println(err)
		return
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

// Prints a list of builds to the user that they can run
func listAction() {
	fmt.Print("#\tNAME\t\tBRANCH\t\tPATH\n--------------------------------------------\n")

	builds := getBuilds(nil)
	for i := 0; i < len(builds); i++ {
		fmt.Printf("%d.\t%s\t\t%s\t\t%s\n", i+1,
			builds[i].Name, builds[i].Branch, builds[i].AbsolutePath)
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
