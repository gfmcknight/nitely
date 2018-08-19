package main

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type argName int

const (
	nameType    argName = 1
	branchType  argName = 2
	remoteType  argName = 3
	dirType     argName = 4
	argsType    argName = 5
	serviceType argName = 6
)

var argNameAliases = map[string]argName{
	"-n":        nameType,
	"--name":    nameType,
	"-b":        branchType,
	"--branch":  branchType,
	"-r":        remoteType,
	"--remote":  remoteType,
	"-d":        dirType,
	"--dir":     dirType,
	"-a":        argsType,
	"--args":    argsType,
	"-s":        serviceType,
	"--service": serviceType,
}

type argSet struct {
	numberedArgs []string
	namedArgs    map[argName]string
}

// Creates an argSet object that can be used for a more flexible
// parsing of arguments
func makeArgSet(args []string) argSet {
	argSetObject := argSet{
		namedArgs:    make(map[argName]string),
		numberedArgs: make([]string, 0),
	}

	for i := 0; i < len(args); i++ {
		if args[i][0] == '-' {
			// Take each argument starting with a - and try to turn it into
			// a named arg. Ignore nonsensical arguments.
			argN := argNameAliases[args[i]]
			if argN == 0 {
				continue
			}

			if i < len(args)-1 {
				argSetObject.namedArgs[argN] = args[i+1]
			} else {
				argSetObject.namedArgs[argN] = " "
			}
		}
	}

	for i := 0; i < len(args); i++ {
		if args[i][0] != '-' && (i == 0 || args[i-1][0] != '-') {
			argSetObject.numberedArgs = append(argSetObject.numberedArgs, args[i])
		}
	}

	return argSetObject
}

func (args argSet) hasArg(i interface{}) bool {
	// TODO: Check if this is actually able to tell the difference
	// between the two
	switch i.(type) {
	case int:
		return len(args.numberedArgs) > i.(int)
	case argName:
		_, exists := args.namedArgs[i.(argName)]
		return exists
	}
	return false
}

func (args argSet) getArg(i interface{}) string {
	// TODO: Check if this is actually able to tell the difference
	// between the two
	switch i.(type) {
	case int:
		return args.numberedArgs[i.(int)]
	case argName:
		return args.namedArgs[i.(argName)]
	}
	return ""
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

// Determine from system properties the amount of
// time to wait when starting services
func sleepTime(build string) time.Duration {
	db := openAndCreateStorage()
	defer db.Close()

	timeStr := getProperty(db, build+".delay")
	if timeStr == nil {
		timeStr = getProperty(db, "delay")
	}
	if timeStr == nil {
		// Our default wait after services is 60 seconds
		return 60
	}

	// Unfortunate string-typing on properties
	timeInt, err := strconv.Atoi(*timeStr)
	if err != nil {
		return 60
	}
	return time.Duration(timeInt) * time.Second
}

// Gets the base directory where nitely is installed
func getStorageBase() string {
	return os.Getenv("NitelyPath")
}
