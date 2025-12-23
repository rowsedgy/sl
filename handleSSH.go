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

func (c *cfg) getTmuxInfo(selectedItem Item, port string) (string, string) {
	session := fmt.Sprintf("ssh-%s", selectedItem.name)
	command := ""

	if selectedItem.pubauth {
		command = fmt.Sprintf("ssh -i %s -o StrictHostKeyChecking=no %s@%s; tmux switch-client -l", selectedItem.key, selectedItem.user, selectedItem.ip)
	}
	if selectedItem.tunnel {
		command = fmt.Sprintf("sshpass -p %s ssh -p %s %s@127.0.0.1 -o HostKeyAlgorithms=+ssh-rsa -o StrictHostKeyChecking=no; tmux switch-client -l", selectedItem.password, port, selectedItem.user)
	} else {
		command = fmt.Sprintf("sshpass -p %s ssh -o StrictHostKeyChecking=no -o PreferredAuthentications=password %s@%s; tmux switch-client -l", selectedItem.password, selectedItem.user, selectedItem.ip)
	}

	return session, command
}

func (c *cfg) spawnTmuxSession(name, command string) {
	exec.Command("tmux", "new-session", "-d", "-s", name).Run()

	if os.Getenv("TMUX") != "" {
		c.runTmuxCommand("tmux", "switch-client", "-t", name)
	} else {
		c.runTmuxCommand("tmux", "attach-session", "-t", name)
	}

	exec.Command("tmux", "send-keys", "-t", name, command, "C-m").Run()
}

func (c *cfg) runTmuxCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (c *cfg) spawnSSHSession(selectedItem Item) error {
	name, command := c.getTmuxInfo(selectedItem, "")
	if selectedItem.pubauth {
		c.spawnTmuxSession(name, command)
	}
	if selectedItem.tunnel {
		resultCh := make(chan tunnelResult)

		go func() {
			port, err := c.startSSHTunnel(selectedItem)
			resultCh <- tunnelResult{port: port, err: err}
		}()

		result := <-resultCh
		if result.err != nil {
			return result.err
		}
		name, command = c.getTmuxInfo(selectedItem, fmt.Sprintf("%d", result.port))
		c.spawnTmuxSession(name, command)

	} else {
		c.spawnTmuxSession(name, command)
	}

	return nil
}

func (c *cfg) startSSHTunnel(selectedItem Item) (int, error) {
	tunnIP := c.connections.TunnelHosts[selectedItem.tunnelHost].IP
	jumpHost := fmt.Sprintf("%s@%s:22", selectedItem.user, tunnIP)
	destHost := fmt.Sprintf("%s:22", selectedItem.ip)

	tunnel := NewSSHTunnel(jumpHost, selectedItem.legacy, selectedItem.password, destHost)

	endpointAddr := "127.0.0.1"
	tunnel.Local = NewEndpoint(endpointAddr)
	// tunnel.Log = log.Default()

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
