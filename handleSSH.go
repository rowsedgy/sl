package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tunnelResult struct {
	port int
	err  error
}

func (c *cfg) spawnSSHSession(selectedItem Item) error {
	var cmd *exec.Cmd

	if selectedItem.pubauth {
		cmd = exec.Command("ssh", "-i", selectedItem.key, "-o", "StrictHostKeyChecking=no", selectedItem.user+"@"+selectedItem.ip)
	}
	if selectedItem.tunnel {
		// freePort, err := getFreePort()
		// if err != nil {
		// 	return err
		// }
		resultCh := make(chan tunnelResult)
		// ready := make(chan struct{})

		// start tunnel in background
		go func() {
			// close(ready)
			// if err := startSSHTunnel(selectedItem); err != nil {
			// 	log.Fatal(err)
			// }
			port, err := c.startSSHTunnel(selectedItem)
			resultCh <- tunnelResult{port: port, err: err}
		}()

		// <-ready
		result := <-resultCh
		if result.err != nil {
			return result.err
		}

		cmd = exec.Command("sshpass", "-p", selectedItem.password, "ssh", "-p", fmt.Sprintf("%d", result.port), fmt.Sprintf("%s@127.0.0.1", selectedItem.user), "-o", "HostKeyAlgorithms=+ssh-rsa", "-o", "StrictHostKeyChecking=no")

	} else {
		cmd = exec.Command("sshpass", "-p", selectedItem.password, "ssh", "-o", "StrictHostKeyChecking=no", "-o", "PreferredAuthentications=password", selectedItem.user+"@"+selectedItem.ip)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

//	func startSSHTunnel(selectedItem Item) (int, error) {
//		jumpHost := fmt.Sprintf("%s@%s:22", selectedItem.user, selectedItem.tunnelHost)
//		destHost := fmt.Sprintf("%s:22", selectedItem.ip)
//		tunnel := NewSSHTunnel(jumpHost, selectedItem.password, destHost)
//		// port := tunnel.Local.Port
//		endpointAddr := fmt.Sprintf("127.0.0.1")
//		tunnel.Local = NewEndpoint(endpointAddr)
//		tunnel.Log = log.Default()
//		port := tunnel.Local.Port
//		fmt.Printf("Tunnel Info:\nTunnelHost: %s\nDestinationHost: %s\nLocalPort: %d\n", jumpHost, destHost, port)
//		return port, tunnel.Start()
//	}
func (c *cfg) startSSHTunnel(selectedItem Item) (int, error) {
	tunnIP := c.connections.TunnelHosts[selectedItem.tunnelHost].IP
	jumpHost := fmt.Sprintf("%s@%s:22", selectedItem.user, tunnIP)
	destHost := fmt.Sprintf("%s:22", selectedItem.ip)

	tunnel := NewSSHTunnel(jumpHost, selectedItem.password, destHost)

	endpointAddr := "127.0.0.1"
	tunnel.Local = NewEndpoint(endpointAddr)
	tunnel.Log = log.Default()

	// start the blocking tunnel in the background
	go func() {
		if err := tunnel.Start(); err != nil {
			log.Printf("tunnel error: %v", err)
		}
	}()

	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(2 * time.Millisecond)
	defer ticker.Stop()

	// check for port until timeout or port != 0 and then return
	for {
		select {
		case <-timeout:
			return 0, fmt.Errorf("timeout waiting for tunnel port")
		case <-ticker.C:
			port := tunnel.Local.Port
			if port != 0 {
				return port, nil
			}
		}
	}
}

func (c *cfg) handleSSHSession(m tea.Model) error {
	model, ok := m.(model)
	if !ok {
		os.Exit(0)
	}

	if model.list.SelectedItem() != nil && model.startSSH {
		selected := model.list.SelectedItem().(Item)
		log.Printf("Connecting to host %s: %s\n", selected.name, selected.ip)

		err := c.spawnSSHSession(selected)
		if err != nil {
			return err
		}
	}

	return nil
}
