package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Subcommands
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("ls", flag.ExitOnError)
	rmCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		fmt.Println("Either of add/ls/remove/update command is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
	case "ls":
		listCmd.Parse(os.Args[2:])
	case "remove":
		rmCmd.Parse(os.Args[2:])
	case "update":
		updateCmd.Parse(os.Args[2:])
	default:
		fmt.Println("No command found")
		os.Exit(1)
	}

}
