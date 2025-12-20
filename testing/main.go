package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	user := "media"
	host := "mediaserver"

	cmd := exec.Command("ssh", user+"@"+host)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
