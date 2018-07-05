package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	//args := os.Args[1:]
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	fmt.Printf(dir)

	if err != nil {
		fmt.Printf("Something went wrong getting the current directory")
		return
	}
}

func getFlags(args []string) map[string]string {
	flags := make(map[string]string)

	for i := 0; i < len(args); i++ {
		if args[i][0] == '-' {
			if i < len(args)-1 {
				flags[args[i]] = args[i+1]
			} else {
				flags[args[i]] = ""
			}
		}
	}

	return flags
}
