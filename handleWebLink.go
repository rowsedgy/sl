package main

import "os/exec"

// implement check for wsl / native linux
// https://gist.github.com/sevkin/9798d67b2cb9d07cb05f89f14ba682f8

func handleWebLink(m model) error {
	url := m.list.SelectedItem().(Item).webip

	cmd := exec.Command("xdg-open", url)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
