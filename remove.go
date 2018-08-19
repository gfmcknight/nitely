package main

import "fmt"

// Removes a build or service from the database
func removeAction(args argSet) {
	if !args.hasArg(1) {
		fmt.Println("Please specify the item to remove from nitely.")
		return
	}

	if args.hasArg(serviceType) {
		deleteServiceInfo(nil, args.getArg(1))
	} else {
		deleteBuildInfo(nil, args.getArg(1))
	}
}
