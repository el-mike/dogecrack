package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	sshPassword := os.Getenv("SSH_PASSWORD")
	sshUser := os.Getenv("SSH_USER")

	hostKeyCb, err := knownhosts.New("~/.ssh/known_hosts")
	if err != nil {
		log.Fatal(err)
		return
	}

	sshHost := "172.20.0.2:22"

	sshConfig := &ssh.ClientConfig{
		User:            sshUser,
		HostKeyCallback: hostKeyCb,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshPassword),
		},
	}

	conn, err := ssh.Dial("tcp", sshHost, sshConfig)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer conn.Close()
}
