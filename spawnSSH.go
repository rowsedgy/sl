package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func spawnSSHSession(user, password, ip string) error {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		return fmt.Errorf("error establishing ssh connection: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("error establishing new session: %v", err)
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	if err := session.Run("bash"); err != nil {
		return fmt.Errorf("error running bash inside session: %v", err)
	}
	return nil
}
