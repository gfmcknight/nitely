package main

import (
	"fmt"
	"os"
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
	case "start":
		startService(nil, argSet.getArg(1))
	case "save":
		db := openAndCreateStorage()
		saveResults(db, getBuildInfo(db, argSet.getArg(1)), argSet.getArg(2))
	}

	fmt.Println()
	fmt.Println(getZen())
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
