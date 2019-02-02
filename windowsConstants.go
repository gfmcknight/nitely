// +build windows

package main

var buildFiles = []string{"nitely.ps1", "nitely.bat", "nitely.py"}
var successFiles = []string{"nitely-success.ps1", "nitely-success.bat", "nitely-success.py"}
var shell = "powershell"
var dependencyFile = "dependencies-nitely.txt"
var resultsFile = "results-nitely.txt"
