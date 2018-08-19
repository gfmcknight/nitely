package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
)

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
