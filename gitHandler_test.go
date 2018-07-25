package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"testing"
)

func fileContiansLine(filepath, filename, line string) bool {
	fmt.Printf("Reading file: %s\n", filename)

	file, err := os.Open(path.Join(filepath, filename))
	defer file.Close()
	if err != nil {
		fmt.Print(err)
		return false
	}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Printf("Next line: %s\n", scanner.Text())
		if scanner.Text() == line {
			return true
		}
	}
	return false
}

func TestInflateCommit(t *testing.T) {
	repoPath := path.Join(os.ExpandEnv("$HOME"), "nitely-test-repo")
	snapshotPath := path.Join(getStorageBase(), "MySnapshot")

	inflateRef("refs/heads/other2", repoPath, "MySnapshot")

	if !fileContiansLine(snapshotPath, "file-a.txt", "ADDITION") {
		t.Error("File a should contain line \"ADDITION\"")
	}

	if !fileContiansLine(snapshotPath, "file-c.txt", "TEST C") {
		t.Error("File a should contain line \"TEST C\"")
	}

	snapshotPath = path.Join(getStorageBase(), "MyOtherSnapshot")

	inflateRef("refs/heads/master", repoPath, "MyOtherSnapshot")

	if !fileContiansLine(snapshotPath, "file-a.txt", "WORKSPACE LINE") {
		t.Error("File a should contain line \"WORKSPACE LINE\"")
	}
}
