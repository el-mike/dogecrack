package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

const CONN_PROTOCOL = "tcp"
const SSH_PORT = 22

// VastClient - Vast.ai SSH connection client. Encapsulates all the operations
// we can perform on Vast.ai instance.
type VastClient struct {
	user      string
	password  string
	ipAddress string
	hostKeyCb ssh.HostKeyCallback

	conn *ssh.Client
}

// NewVastClient  - returns new VastClient instance. Does NOT start up a connection.
func NewVastClient(user, password, sshDirPath, ipAddress string) (*VastClient, error) {
	hostKeyCb, err := knownhosts.New(sshDirPath + "/known_hosts")
	if err != nil {
		return nil, err
	}

	client := &VastClient{
		user:      user,
		password:  password,
		hostKeyCb: hostKeyCb,
		ipAddress: ipAddress,
	}

	return client, nil
}

// Connect - starts a SSH connection.
func (vs *VastClient) Connect() error {
	config := &ssh.ClientConfig{
		User:            vs.user,
		HostKeyCallback: vs.hostKeyCb,
		Auth: []ssh.AuthMethod{
			ssh.Password(vs.password),
		},
	}

	host := vs.ipAddress + ":" + fmt.Sprint(SSH_PORT)

	conn, err := ssh.Dial(CONN_PROTOCOL, host, config)
	if err != nil {
		return err
	}

	vs.conn = conn

	return nil
}

// Close - closes the connection.
func (vs *VastClient) Close() error {
	return vs.conn.Close()
}

// GetUser - prints current user to stdout.
func (vs *VastClient) GetUser() error {
	return vs.run("whoami", 10)
}

func (vs *VastClient) run(cmd string, timeout int) error {
	session, err := vs.conn.NewSession()
	if err != nil {
		return err
	}

	defer session.Close()

	session.Stderr = os.Stderr

	sessionStdOut, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	sessionStdIn, err := session.StdinPipe()
	if err != nil {
		return err
	}

	go fmt.Fprintf(sessionStdIn, cmd+"\n")
	err = session.Shell()
	if err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(sessionStdOut)

		for scanner.Scan() {
			output := scanner.Text()
			fmt.Printf("%s\n", output)
		}
	}()

	time.Sleep(time.Duration(timeout) * time.Second)

	return nil
}
