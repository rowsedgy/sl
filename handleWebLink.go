package main

import (
	"os/exec"
	"runtime"
	"strings"
)

// implement check for wsl / native linux
// https://gist.github.com/sevkin/9798d67b2cb9d07cb05f89f14ba682f8

func handleWebLink(m model) error {
	url := m.list.SelectedItem().(Item).webip

	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	default:
		if isWSL() {
			cmd = "cmd.exe"
			args = []string{"/c", "start", url}
		} else {
			cmd = "xdg-open"
			args = []string{url}
		}
	}

	return exec.Command(cmd, args...).Start()
}

func isWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}
