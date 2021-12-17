package main

import (
	"fmt"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

const CONN_PROTOCOL = "tcp"
const SSH_PORT = 22

type VastClient struct {
	conn *ssh.Client
}

func NewVastClient(user, password, sshDirPath, ipAddress string) (*VastClient, error) {
	client := &VastClient{}

	hostKeyCb, err := knownhosts.New(sshDirPath + "/known_hosts")
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: hostKeyCb,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}

	host := ipAddress + ":" + fmt.Sprint(SSH_PORT)

	conn, err := ssh.Dial(CONN_PROTOCOL, host, config)
	if err != nil {
		return nil, err
	}

	client.conn = conn

	return client, nil
}
