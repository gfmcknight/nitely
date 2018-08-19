package main

import "fmt"

// Prints a list of builds to the user that they can run
func listAction(args argSet) {
	if !args.hasArg(1) {
		fmt.Println("Please specify whether to list builds or services.")
		return
	}

	switch args.getArg(1) {

	case "builds":
		fmt.Print("#\tNAME\t\tBRANCH\t\tPATH\t\tREMOTE\n--------------------------------------------\n")
		builds := getBuilds(nil)
		for i := 0; i < len(builds); i++ {
			fmt.Printf("%d.\t%s\t%s\t%s\t%s\n", i+1,
				builds[i].Name, builds[i].Branch, builds[i].AbsolutePath, builds[i].Remote)
		}

	case "services":
		fmt.Print("#\tNAME\t\tARGS\t\tPATH\n--------------------------------------------\n")
		services := getServices(nil)
		for i := 0; i < len(services); i++ {
			fmt.Printf("%d.\t%s\t%s\t%s\n", i+1,
				services[i].Name, services[i].Args, services[i].AbsolutePath)
		}
	}
}
