package main

import "testing"

func TestArgParser(t *testing.T) {
	args := makeArgSet([]string{"arg1", "arg2", "-n", "nameArg", "--dir", "dirArg", "arg3"})

	if !args.hasArg(2) {
		t.Fail()
	}
	if args.hasArg(3) {
		t.Fail()
	}

	if args.hasArg(branchType) {
		t.Fail()
	}
	if !args.hasArg(dirType) {
		t.Fail()
	}

	if args.getArg(0) != "arg1" {
		t.Fail()
	}
	if args.getArg(1) != "arg2" {
		t.Fail()
	}
	if args.getArg(2) != "arg3" {
		t.Fail()
	}
	if args.getArg(nameType) != "nameArg" {
		t.Fail()
	}
	if args.getArg(dirType) != "dirArg" {
		t.Fail()
	}
}
