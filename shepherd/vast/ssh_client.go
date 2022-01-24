package vast

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

const CONN_PROTOCOL = "tcp"
const CONTAINER_PITBULL_PATH = "/app"

// VastSSHClient - Vast.ai SSH connection client. Encapsulates all the operations
// we can perform on Vast.ai instance.
type VastSSHClient struct {
	user      string
	password  string
	hostKeyCb ssh.HostKeyCallback

	ipAddress string
	port      int

	conn *ssh.Client
}

// NewVastSSHClient  - returns new VastClient instance. Does NOT start up a connection.
func NewVastSSHClient(user, password, sshDirPath, ipAddress string, port int) (*VastSSHClient, error) {
	knownHostsPath := sshDirPath + "/known_hosts"

	hostKeyCb, err := knownhosts.New(knownHostsPath)
	if err != nil {
		return nil, err
	}

	client := &VastSSHClient{
		user:      user,
		password:  password,
		hostKeyCb: hostKeyCb,
		ipAddress: ipAddress,
		port:      port,
	}

	return client, nil
}

// RunPitbull - runs Pitbull process for given fileUrl and walletString.
func (vs *VastSSHClient) RunPitbull(fileUrl, walletString string) (string, error) {
	return vs.run("pitbull run -f " + fileUrl + " -w " + walletString)
}

// GetPitbullStatus - runs Pitbull's status command and returns the output.
func (vs *VastSSHClient) GetPitbullStatus() (string, error) {
	return vs.run("pitbull status")
}

// GetPitbullStatus - runs Pitbull's progress command and returns the output.
func (vs *VastSSHClient) GetPitbullProgress() (string, error) {
	return vs.run("pitbull progress")
}

// GetPitbullOutput - runs Pitbull's output command and returns the output.
func (vs *VastSSHClient) GetPitbullOutput() (string, error) {
	return vs.run("pitbull output")
}

// Connect - starts a SSH connection.
func (vs *VastSSHClient) connect() error {
	config := &ssh.ClientConfig{
		User:            vs.user,
		HostKeyCallback: vs.hostKeyCb,
		Auth: []ssh.AuthMethod{
			ssh.Password(vs.password),
		},
	}

	host := vs.ipAddress + ":" + fmt.Sprint(vs.port)

	conn, err := ssh.Dial(CONN_PROTOCOL, host, config)
	if err != nil {
		return err
	}

	vs.conn = conn

	return nil
}

// Close - closes the connection.
func (vs *VastSSHClient) close() error {
	return vs.conn.Close()
}

func (vs *VastSSHClient) run(cmd string) (string, error) {
	if err := vs.connect(); err != nil {
		return "", err
	}
	defer vs.close()

	session, err := vs.conn.NewSession()
	if err != nil {
		return "", err
	}

	defer session.Close()

	session.Stderr = os.Stderr

	sessionStdOut, err := session.StdoutPipe()
	if err != nil {
		return "", err
	}

	err = session.Run(cmd)
	if err != nil {
		// If the error was an ExitError, we probably want to proceed, as some
		// Pitbull scripts return custom exit codes.
		if err, ok := err.(*ssh.ExitError); !ok {
			return "", err
		}
	}

	scanner := bufio.NewScanner(sessionStdOut)

	output := ""

	for scanner.Scan() {
		line := scanner.Text()
		output += line
	}

	return output, nil
}
