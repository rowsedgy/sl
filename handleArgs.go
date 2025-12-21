package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func (c *cfg) handleArgs(args []string) {
	if args[0] == "ls" && len(args) < 2 {
		err := c.listEntries()
		if err != nil {
			fmt.Println("ERROR -", err)
		}
		os.Exit(0)
	}
	if args[0] == "--help" || len(args) < 2 {
		c.printHelp()
	}

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addName := addCmd.String("name", "None", "HostName to add")
	addIP := addCmd.String("ip", "None", "Host SSH IP")
	addWebIP := addCmd.String("webip", "None", "Web Interface IP")
	addUser := addCmd.String("user", "None", "SSH User")
	addPassword := addCmd.String("password", "None", "SSH Password")
	addPubauth := addCmd.Bool("pubauth", false, "PubKey authentication")
	addPubKey := addCmd.String("key", "None", "Key file path")

	removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	removeName := removeCmd.String("name", "Nil", "HostName to remove")

	switch args[0] {
	case "add":
		addCmd.Parse(args[1:])
		err := c.addEntry(*addName, *addIP, *addWebIP, *addUser, *addPassword, *addPubKey, *addPubauth)
		if err != nil {
			fmt.Println("ERROR -", err)
		}
		os.Exit(0)
	case "remove":
		removeCmd.Parse(args[1:])
		err := c.removeEntry(*removeName)
		if err != nil {
			fmt.Println("ERROR -", err)
		}
		os.Exit(0)
	default:
		c.printHelp()
	}
}

func (c *cfg) addEntry(name, ip, webip, user, password, key string, pubauth bool) error {
	bytes, err := os.ReadFile(c.filepath)
	if err != nil {
		return err
	}

	var connections []connection

	if len(bytes) > 3 {
		err = json.Unmarshal(bytes, &connections)
		if err != nil {
			return fmt.Errorf("Error unmarshaling json in addentry: %v", err)
		}

		// check duplicates
		for _, c := range connections {
			if c.Name == name {
				return fmt.Errorf("Entry %s already exists", name)
			}
		}
	}

	var newEntry connection
	newEntry.Name = name
	newEntry.Data.IP = ip
	newEntry.Data.WebIP = webip
	newEntry.Data.User = user
	newEntry.Data.Password = password
	newEntry.Data.Pubauth = pubauth
	newEntry.Data.Key = key

	connections = append(connections, newEntry)

	updatedBytes, err := json.MarshalIndent(connections, "", "\t")
	if err != nil {
		return fmt.Errorf("Error marshaling new json: %v", err)
	}

	err = os.WriteFile(c.filepath, updatedBytes, 0644)
	if err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	return nil
}

func (c *cfg) removeEntry(name string) error {
	bytes, err := os.ReadFile(c.filepath)
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

	err = os.WriteFile(c.filepath, updatedBytes, 0644)
	if err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	return nil
}

func (c *cfg) listEntries() error {
	bytes, err := os.ReadFile(c.filepath)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}

	var connections []connection

	if len(bytes) < 3 {
		return fmt.Errorf("No entries in %s", c.filepath)
	}

	err = json.Unmarshal(bytes, &connections)
	if err != nil {
		return fmt.Errorf("Error unmarshaling json: %v", err)
	}

	for _, c := range connections {
		fmt.Printf("* Name: %s\n - IP: %s\n - WebIp: %s\n", c.Name, c.Data.IP, c.Data.WebIP)
	}

	return nil
}

func (c *cfg) printHelp() {
	cmd := os.Args[0]
	fmt.Println("Usage:")
	fmt.Printf("* %s\n\tLaunches list in interactive mode\n", cmd)
	fmt.Printf("* %s ls\n\tLists all available host names\n", cmd)
	fmt.Printf("* %s add --name=<hostname> --ip=<ip address> --webip=<web interface addr> --user=<ssh username> --password=<ssh password>\n\t Adds an entry to the list for the provided host\n", cmd)
	fmt.Printf("* %s remove --name<hostname>\n\tRemoves provided host from entries\n", cmd)
	fmt.Printf("\n\nConnection file is located at %s\n", c.filepath)

	os.Exit(0)
}
