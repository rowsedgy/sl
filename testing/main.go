package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// add command
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addName := addCmd.String("name", "None", "HostName")
	addIP := addCmd.String("ip", "None", "SSH IP address")
	addWebIP := addCmd.String("webip", "None", "Web interface address")

	// remove command
	removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	removeName := removeCmd.String("name", "nil", "HostName to remove")

	if len(os.Args) < 2 {
		fmt.Println("not enough args")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		fmt.Println("Adding host:")
		fmt.Println("	name:", *addName)
		fmt.Println("	ip:", *addIP)
		fmt.Println("	webip:", *addWebIP)
	case "remove":
		removeCmd.Parse(os.Args[2:])
		fmt.Println("Removing host:")
		fmt.Println("	name:", *removeName)
	case "ls":
		fmt.Println("list command executed")
	default:
		fmt.Println("No subcommand provided")
		os.Exit(1)
	}
}
