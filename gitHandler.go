package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
)

type priorStateInfo struct {
	branchChanged  bool
	prevBranch     string
	changesStashed bool
}

type objectInfo struct {
	tree bool
	id   string
	name string
}

func writeObjectInfo(object, dest string) {
	cmd := exec.Command("git", "cat-file", "-p", object)
	file, err := os.Open(dest)
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	defer file.Close()
	cmd.Stdout = file

	err = cmd.Run()

	if err != nil {
		fmt.Print(err)
		panic(err)
	}
}

func getObjectInfo(object string) []byte {
	cmd := exec.Command("git", "cat-file", "-p", object)
	output, err := cmd.Output()
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	return output
}

func inflateBlob(currentPath, object, name string) {
	writeObjectInfo(object, path.Join(currentPath, name))
}

func inflateTree(currentPath, object, name string) {
	newPath := path.Join(currentPath, name)
	os.MkdirAll(newPath, os.ModeDir)
	scanner := bufio.NewScanner(bytes.NewReader(getObjectInfo(object)))

	objects := make([]objectInfo, 0)

	for scanner.Scan() {
		if scanner.Text() != "blob" && scanner.Text() != "tree" {
			continue
		}

		newObject := objectInfo{}
		newObject.tree = (scanner.Text() == "tree")
		scanner.Scan()
		newObject.id = scanner.Text()
		scanner.Scan()
		newObject.name = scanner.Text()
		objects = append(objects, newObject)
	}

	for _, newObject := range objects {
		if newObject.tree {
			inflateTree(newPath, newObject.id, newObject.name)
		} else {
			inflateBlob(newPath, newObject.id, newObject.name)
		}
	}
}

func inflateCommit(ref, src, snapshotName string) {
	os.Chdir(src)
	dir := getStorageBase()
	dest := path.Join(dir, snapshotName)

	currentRef := string(getObjectInfo("HEAD"))
	if currentRef == ref {
		// If the head is already at a given branch, then we
		// will want to
		copyDir(src, dest)
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(getObjectInfo(ref)))
	treeID := ""
	for scanner.Scan() {
		if scanner.Text() == "tree" {
			scanner.Scan()
			treeID = scanner.Text()
			break
		}
	}
	if treeID == "" {
		return
	}

	inflateTree(dir, treeID, snapshotName)

}
