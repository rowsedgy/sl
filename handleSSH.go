package main

import (
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

func spawnSSHSession(user, password, ip string) error {
	cmd := exec.Command("sshpass", "-p", password, "ssh", "-o", "StrictHostKeyChecking=no", "-o", "PreferredAuthentications=password", user+"@"+ip)

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

		err := spawnSSHSession(selected.user, selected.password, selected.ip)
		if err != nil {
			return err
		}
	}

	return nil
}
