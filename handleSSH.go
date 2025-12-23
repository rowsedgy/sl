package main

import (
	"log"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

func spawnSSHSession(selectedItem Item) error {
	var cmd *exec.Cmd

	if selectedItem.pubauth {
		cmd = exec.Command("ssh", "-i", selectedItem.key, "-o", "StrictHostKeyChecking=no", selectedItem.user+"@"+selectedItem.ip)
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

func handleSSHSession(m tea.Model) error {
	model, ok := m.(model)
	if !ok {
		os.Exit(0)
	}

	if model.list.SelectedItem() != nil && model.startSSH {
		selected := model.list.SelectedItem().(Item)
		log.Printf("Connecting to host %s: %s\n", selected.name, selected.ip)

		err := spawnSSHSession(selected)
		if err != nil {
			return err
		}
	}

	return nil
}
