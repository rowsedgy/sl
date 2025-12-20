package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

// const testFile = "connections2.json"

func handleArgs(args []string) {
	if args[0] == "ls" && len(args) < 2 {
		err := listEntries(filepath)
		if err != nil {
			fmt.Println("ERROR -", err)
		}
		os.Exit(0)
	}
	if args[0] == "--help" || len(args) < 2 {
		printHelp()
	}

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addName := addCmd.String("name", "None", "HostName to add")
	addIP := addCmd.String("ip", "None", "Host SSH IP")
	addWebIP := addCmd.String("webip", "None", "Web Interface IP")
	addUser := addCmd.String("user", "None", "SSH User")
	addPassword := addCmd.String("password", "None", "SSH Password")

	removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	removeName := removeCmd.String("name", "Nil", "HostName to remove")

	switch args[0] {
	case "add":
		addCmd.Parse(args[1:])
		err := addEntry(filepath, *addName, *addIP, *addWebIP, *addUser, *addPassword)
		if err != nil {
			fmt.Println("ERROR -", err)
		}
		os.Exit(0)
	case "remove":
		removeCmd.Parse(args[1:])
		err := removeEntry(filepath, *removeName)
		if err != nil {
			fmt.Println("ERROR -", err)
		}
		os.Exit(0)
	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Println("printing help")
	os.Exit(0)
}

func addEntry(file, name, ip, webip, user, password string) error {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	var connections []connection

	err = json.Unmarshal(bytes, &connections)
	if err != nil {
		return fmt.Errorf("Error unmarshaling json: %v", err)
	}

	// check duplicates
	for _, c := range connections {
		if c.Name == name {
			return fmt.Errorf("Entry %s already exists", name)
		}
	}

	var newEntry connection
	newEntry.Name = name
	newEntry.Data.IP = ip
	newEntry.Data.WebIP = webip
	newEntry.Data.User = user
	newEntry.Data.Password = password

	connections = append(connections, newEntry)

	updatedBytes, err := json.MarshalIndent(connections, "", "\t")
	if err != nil {
		return fmt.Errorf("Error marshaling new json: %v", err)
	}

	err = os.WriteFile(file, updatedBytes, 0644)
	if err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	return nil
}

func removeEntry(file, name string) error {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}

	var connections []connection

	err = json.Unmarshal(bytes, &connections)
	if err != nil {
		return fmt.Errorf("Error unmarshaling json: %v", err)
	}

	for i, c := range connections {
		if c.Name == name {
			connections = append(connections[:i], connections[i+1:]...)
		}
	}

	updatedBytes, err := json.MarshalIndent(connections, "", "\t")
	if err != nil {
		return fmt.Errorf("Error marshaling new json: %v", err)
	}

	err = os.WriteFile(file, updatedBytes, 0644)
	if err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	return nil
}

func listEntries(file string) error {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}

	var connections []connection

	err = json.Unmarshal(bytes, &connections)
	if err != nil {
		return fmt.Errorf("Error unmarshaling json: %v", err)
	}

	for _, c := range connections {
		fmt.Printf("- Name: %s\n\t- IP: %s\n\t- WebIp: %s\n", c.Name, c.Data.IP, c.Data.WebIP)
	}

	return nil
}
