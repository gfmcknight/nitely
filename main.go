package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	fmt.Printf(dir)

	if err != nil {
		fmt.Printf("Something went wrong getting the current directory")
		return
	}

	if len(args) < 1 {
		fmt.Printf("USAGE")
		return
	}

	switch args[0] {
	case "add":
		addAction(dir, getFlags(args))
	case "remove":
		removeAction(args[1], getFlags(args))
	case "build":
		buildAction(args[1], getFlags(args))
	}
}

func addAction(dir string, flags map[string]string) {
	if newDir, exists := flags["-d"]; exists {
		dir = filepath.Join(dir, newDir)
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
	fmt.Printf("Added build %s on branch %s\nWith path %s", toAdd.Name, toAdd.Branch, toAdd.AbsolutePath)
}

func removeAction(name string, flags map[string]string) {

}

// Runs a build of the given repository
func buildAction(name string, flags map[string]string) {
	/*buildInfo := getBuildInfo(nil, name)
	priorState := prepareRepository(buildInfo.AbsolutePath, buildInfo.Branch)
	defer restoreRepo(buildInfo.AbsolutePath, priorState)

	os.Chdir(buildInfo.AbsolutePath)
	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Print(err)
		return
	}

	for _, filename := range buildFiles {
		if containsFile(filename, files) {
			// TODO: Run here
			return
		}
	}*/

	fmt.Printf("No suitable nitely file found to build")
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

			fmt.Printf("Added flag %s with value %s", args[i], flags[args[i]])
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
