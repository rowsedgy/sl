package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Endpoint struct {
	Host string
	Port int
	User string
}

func NewEndpoint(s string) *Endpoint {
	endpoint := &Endpoint{
		Host: s,
	}

	if parts := strings.Split(endpoint.Host, "@"); len(parts) > 1 {
		endpoint.User = parts[0]
		endpoint.Host = parts[1]
	}

	if parts := strings.Split(endpoint.Host, ":"); len(parts) > 1 {
		endpoint.Host = parts[0]
		endpoint.Port, _ = strconv.Atoi(parts[1])
	}

	return endpoint
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

type SSHTunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint
	Config *ssh.ClientConfig
	Log    *log.Logger
}

func (t *SSHTunnel) logf(format string, args ...any) {
	if t.Log != nil {
		t.Log.Printf(format, args...)
	}
}

func (t *SSHTunnel) Start() error {
	listener, err := net.Listen("tcp", t.Local.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	t.Local.Port = listener.Addr().(*net.TCPAddr).Port
	t.logf("listening on %s", t.Local.String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		t.logf("accepted connection")
		go t.forward(conn)
	}
}

func (t *SSHTunnel) forward(localConn net.Conn) {
	defer localConn.Close()

	serverConn, err := ssh.Dial("tcp", t.Server.String(), t.Config)
	if err != nil {
		t.logf("server dial error: %v", err)
		return
	}
	defer serverConn.Close()

	t.logf("connected to %s (1/2)", t.Server.String())

	remoteConn, err := serverConn.Dial("tcp", t.Remote.String())
	if err != nil {
		t.logf("remote dial error: %v", err)
		return
	}
	defer remoteConn.Close()

	t.logf("connected to %s (2/2)", t.Remote.String())

	go io.Copy(remoteConn, localConn)
	io.Copy(localConn, remoteConn)
}

func PrivateKeyFile(path string) (ssh.AuthMethod, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(key), nil
}

func NewSSHTunnel(tunnel string, legacy bool, password, destination string) *SSHTunnel {
	// if port == 0 a random port will be chosen
	localEndpoint := NewEndpoint("localhost:0")
	server := NewEndpoint(tunnel)
	if server.Port == 0 {
		server.Port = 22
	}
	authType := []ssh.AuthMethod{ssh.Password(password)}

	if legacy {
		authType = []ssh.AuthMethod{
			ssh.KeyboardInteractive(
				func(user, instruction string, questions []string, echos []bool) ([]string, error) {
					answers := make([]string, len(questions))
					for i := range answers {
						answers[i] = password
					}
					return answers, nil
				},
			),
		}
	}
	sshTunnel := &SSHTunnel{
		Config: &ssh.ClientConfig{
			User: server.User,
			Auth: authType,
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				// Always accept key.
				return nil
			},
		},
		Local:  localEndpoint,
		Server: server,
		Remote: NewEndpoint(destination),
	}
	return sshTunnel
}
