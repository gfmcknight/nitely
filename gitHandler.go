package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
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
	file, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer file.Close()
	cmd.Stdout = file

	err = cmd.Run()

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func getHeadRef() string {
	cmd := exec.Command("git", "symbolic-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return strings.Replace(string(output), "\n", "", -1)
}

func getObjectInfo(object string) []byte {
	cmd := exec.Command("git", "cat-file", "-p", object)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
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
	scanner.Split(bufio.ScanWords)

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

func inflateRef(ref, src, snapshotName string) {
	src = path.Clean(src)

	os.Chdir(src)
	dir := getStorageBase()
	dest := path.Join(dir, snapshotName)

	currentRef := getHeadRef()
	if ref == "" || currentRef == ref {
		// If the head is already at a given branch, then we
		// will want to copy it to test the changes they haven't
		// yet committed; if the user didn't specify a
		// branch then we want to use the working tree
		copyDir(src, dest)
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(getObjectInfo(ref)))
	scanner.Split(bufio.ScanWords)
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
