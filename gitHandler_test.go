package main

import "testing"

func TestInflateCommit(t *testing.T) {
	inflateCommit("refs/heads/master", "", "MySnapshot")
}
