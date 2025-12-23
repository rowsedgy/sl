package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func (c *cfg) handleArgs(args []string) {
	if len(args) < 2 {
		if args[0] == "ls" {
			err := c.listEntries()
			if err != nil {
				fmt.Println("ERROR -", err)
			}
			os.Exit(0)
		}
		if args[0] == "lstun" {
			err := c.listTunnels()
			if err != nil {
				fmt.Println("ERROR -", err)
			}
			os.Exit(0)
		}
	}

	addHostCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addName := addHostCmd.String("name", "None", "HostName to add")
	addIP := addHostCmd.String("ip", "None", "Host SSH IP")
	addWebIP := addHostCmd.String("webip", "None", "Web Interface IP")
	addUser := addHostCmd.String("user", "None", "SSH User")
	addPassword := addHostCmd.String("password", "None", "SSH Password")
	addPubauth := addHostCmd.Bool("pubauth", false, "PubKey authentication")
	addPubKey := addHostCmd.String("key", "None", "Key file path")
	addTunnel := addHostCmd.Bool("tunnel", false, "Connect through tunnel")
	addTunnelHost := addHostCmd.String("tunnelhost", "None", "Tunnel Host")

	addTunnelCmd := flag.NewFlagSet("addtun", flag.ExitOnError)
	addTunnelName := addTunnelCmd.String("name", "None", "Tunnel hostname to add")
	addTunnelIP := addTunnelCmd.String("ip", "None", "Tunnel host IP")
	addTunnelUser := addTunnelCmd.String("user", "None", "Tunnel ssh user")
	addTunnelPassword := addTunnelCmd.String("password", "None", "Tunnel ssh password")

	removeHostCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	removeName := removeHostCmd.String("name", "Nil", "HostName to remove")

	removeTunnelCmd := flag.NewFlagSet("removetun", flag.ExitOnError)
	removeTunnelName := removeTunnelCmd.String("name", "Nil", "Tunnel Hostname to remove")

	switch args[0] {
	case "add":
		addHostCmd.Parse(args[1:])
		err := c.addEntry(*addName, *addIP, *addWebIP, *addUser, *addPassword, *addPubKey, *addTunnelHost, *addPubauth, *addTunnel)
		if err != nil {
			fmt.Println("ERROR -", err)
		}
		os.Exit(0)
	case "addtun":
		addTunnelCmd.Parse(args[1:])
		err := c.addTunnel(*addTunnelName, *addTunnelUser, *addTunnelPassword, *addTunnelIP)
		if err != nil {
			fmt.Println("ERROR -", err)
		}
		os.Exit(0)
	case "remove":
		removeHostCmd.Parse(args[1:])
		err := c.removeEntry(*removeName)
		if err != nil {
			fmt.Println("ERROR -", err)
		}
		os.Exit(0)
	case "removetun":
		removeTunnelCmd.Parse(args[1:])
		err := c.removeTunnel(*removeTunnelName)
		if err != nil {
			fmt.Println("ERROR -", err)
		}
		os.Exit(0)
	default:
		c.printHelp()
	}
}

func (c *cfg) addEntry(name, ip, webip, user, password, key, tunnelhost string, pubauth, tunnel bool) error {
	bytes, err := os.ReadFile(c.filepath)
	if err != nil {
		return err
	}

	var conns connections

	if len(bytes) > 3 {
		err = json.Unmarshal(bytes, &conns)
		if err != nil {
			return fmt.Errorf("Error unmarshaling json in addentry: %v", err)
		}
	}

	// check duplicates
	if _, ok := conns.Hosts[name]; ok {
		return fmt.Errorf("Entry %s already exists", name)
	}

	newEntry := Host{}
	newEntry.IP = ip
	newEntry.WebIP = webip
	newEntry.User = user
	newEntry.Password = password
	newEntry.Pubauth = pubauth
	newEntry.Key = key
	newEntry.TunnelHost = tunnelhost
	newEntry.Tunnel = tunnel

	conns.Hosts[name] = newEntry

	updatedBytes, err := json.MarshalIndent(conns, "", "\t")
	if err != nil {
		return fmt.Errorf("Error marshaling new json: %v", err)
	}

	err = os.WriteFile(c.filepath, updatedBytes, 0o644)
	if err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	return nil
}

func (c *cfg) addTunnel(name, user, password, ip string) error {
	bytes, err := os.ReadFile(c.filepath)
	if err != nil {
		return err
	}

	var conns connections

	if len(bytes) > 3 {
		err = json.Unmarshal(bytes, &conns)
		if err != nil {
			return fmt.Errorf("Error unmarshaling json in addentry: %v", err)
		}
	}

	if _, ok := conns.TunnelHosts[name]; ok {
		return fmt.Errorf("Tunnel entry %s already exists", name)
	}

	newTunnelEntry := TunnelHost{}
	newTunnelEntry.IP = ip
	newTunnelEntry.User = user
	newTunnelEntry.Password = password
	conns.TunnelHosts[name] = newTunnelEntry

	updatedBytes, err := json.MarshalIndent(conns, "", "\t")
	if err != nil {
		return fmt.Errorf("Error marshaling new json: %v", err)
	}

	err = os.WriteFile(c.filepath, updatedBytes, 0o644)
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

	var conns connections

	err = json.Unmarshal(bytes, &conns)
	if err != nil {
		return fmt.Errorf("Error unmarshaling json: %v", err)
	}

	if _, ok := conns.Hosts[name]; ok {
		delete(conns.Hosts, name)
		updatedBytes, err := json.MarshalIndent(conns, "", "\t")
		if err != nil {
			return fmt.Errorf("Error marshaling new json: %v", err)
		}

		err = os.WriteFile(c.filepath, updatedBytes, 0o644)
		if err != nil {
			return fmt.Errorf("Error writing to file: %v", err)
		}

		return nil
	}
	return fmt.Errorf("Entry %s not found", name)
}

func (c *cfg) removeTunnel(name string) error {
	bytes, err := os.ReadFile(c.filepath)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}

	var conns connections

	err = json.Unmarshal(bytes, &conns)
	if err != nil {
		return fmt.Errorf("Error unmarshaling json: %v", err)
	}

	if _, ok := conns.TunnelHosts[name]; ok {
		delete(conns.TunnelHosts, name)
		updatedBytes, err := json.MarshalIndent(conns, "", "\t")
		if err != nil {
			return fmt.Errorf("Error marshaling new json: %v", err)
		}

		err = os.WriteFile(c.filepath, updatedBytes, 0o644)
		if err != nil {
			return fmt.Errorf("Error writing to file: %v", err)
		}

		return nil
	}
	return fmt.Errorf("Tunnel %s not found", name)
}

func (c *cfg) listEntries() error {
	bytes, err := os.ReadFile(c.filepath)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}

	var conns connections

	if len(bytes) < 3 {
		return fmt.Errorf("No entries in %s", c.filepath)
	}

	err = json.Unmarshal(bytes, &conns)
	if err != nil {
		return fmt.Errorf("Error unmarshaling json: %v", err)
	}

	// fmt.Println(c.connections.Hosts)
	for name, data := range conns.Hosts {
		fmt.Printf("* Name: %s\n - IP: %s\n - WebIp: %s\n - PubKey: %v\n - KeyPath: %s\n - Tunnel: %v\n - TunnelHost: %s\n", name, data.IP, data.WebIP, data.Pubauth, data.Key, data.Tunnel, data.TunnelHost)
	}
	// for _, c := range connections {
	// 	fmt.Printf("* Name: %s\n - IP: %s\n - WebIp: %s\n - PubKey: %v\n - KeyPath: %s\n", c.Name, c.Data.IP, c.Data.WebIP, c.Data.Pubauth, c.Data.Key)
	// }

	return nil
}

func (c *cfg) listTunnels() error {
	bytes, err := os.ReadFile(c.filepath)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}

	var conns connections

	if len(bytes) < 3 {
		return fmt.Errorf("No entries in %s", c.filepath)
	}

	err = json.Unmarshal(bytes, &conns)
	if err != nil {
		return fmt.Errorf("Error unmarshaling json: %v", err)
	}

	for name, data := range conns.TunnelHosts {
		fmt.Printf("* Name: %s\n - IP: %s\n", name, data.IP)
	}
	return nil
}

func (c *cfg) printHelp() {
	cmd := os.Args[0]
	fmt.Println("Usage:")
	fmt.Printf("* %s\n\tLaunches list in interactive mode\n", cmd)
	fmt.Printf("* %s ls\n\tLists all available host names\n", cmd)
	fmt.Printf("* %s add --name=<hostname> --ip=<ip address> --webip=<web interface addr> --user=<ssh username> --pubkey=<true/false> --key=<pubkey path> --password=<ssh password>\n\t Adds an entry to the list for the provided host\n", cmd)
	fmt.Printf("* %s remove --name<hostname>\n\tRemoves provided host from entries\n", cmd)
	fmt.Printf("\n\nConnection file is located at %s\n", c.filepath)

	os.Exit(0)
}
