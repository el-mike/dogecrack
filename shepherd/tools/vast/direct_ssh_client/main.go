package main

import (
	"bufio"
	"fmt"
	"github.com/el-mike/dogecrack/shepherd/internal/config"
	"golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	rootPath, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}

	appConfig, err := config.NewAppConfig(rootPath)
	if err != nil {
		panic(err)
	}

	instanceIp := os.Getenv("DIRECT_SSH_CLIENT_INSTANCE_IP")
	instancePort := os.Getenv("DIRECT_SSH_CLIENT_INSTANCE_PORT")
	//instanceIp, _ := vast.GetFakeVastIp(rootPath, 1)
	//instancePort := 22

	sshPrivateKey := strings.ReplaceAll(appConfig.SSHPrivateKey, `\n`, "\n")
	privateKeyRaw := []byte(sshPrivateKey)

	signer, err := ssh.ParsePrivateKey(privateKeyRaw)
	if err != nil {
		panic(err)
	}

	config := &ssh.ClientConfig{
		User:            appConfig.SSHUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
			ssh.Password(appConfig.SSHPassword),
		},
	}

	host := instanceIp + ":" + fmt.Sprint(instancePort)

	conn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.Stderr = os.Stderr

	sessionStdOut, err := session.StdoutPipe()
	if err != nil {
		panic(err)
	}

	err = session.Run("pitbull status")
	if err != nil {
		// If the error was an ExitError, we probably want to proceed, as some
		// Pitbull scripts return custom exit codes.
		if err, ok := err.(*ssh.ExitError); !ok {
			panic(err)
		}
	}

	scanner := bufio.NewScanner(sessionStdOut)

	output := ""

	for scanner.Scan() {
		line := scanner.Text()
		output += line + "\n"
	}

	fmt.Println(output)
}
