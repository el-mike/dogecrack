package vast

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

const CONN_PROTOCOL = "tcp"

// VastSSHClient - Vast.ai SSH connection client. Encapsulates all the operations
// we can perform on Vast.ai instance.
type VastSSHClient struct {
	user      string
	password  string
	signer    ssh.Signer
	hostKeyCb ssh.HostKeyCallback

	ipAddress string
	port      int

	conn *ssh.Client
}

// NewVastSSHClient - returns new VastSSHClient instance. Does NOT start up a connection.
func NewVastSSHClient(user, password, sshDirPath, sshPrivateKey, ipAddress string, port int) (*VastSSHClient, error) {
	// We store sshPrivateKey as string with '\n', but since it's coming from .env file, it treats
	// '\n' characters as normal signs. Therefore, we use replace to insert actual new lines.
	sshPrivateKey = strings.ReplaceAll(sshPrivateKey, `\n`, "\n")
	privateKeyRaw := []byte(sshPrivateKey)

	signer, err := ssh.ParsePrivateKey(privateKeyRaw)
	if err != nil {
		return nil, err
	}

	client := &VastSSHClient{
		user:     user,
		password: password,
		signer:   signer,
		// TODO: Consider handling host key verification - it's partially implemented below.
		hostKeyCb: ssh.InsecureIgnoreHostKey(),
		ipAddress: ipAddress,
		port:      port,
	}

	return client, nil
}

// RunPitbullForPasslist - runs Pitbull process for given passlistUrl and walletString.
func (vs *VastSSHClient) RunPitbullForPasslist(passlistUrl, walletString string) (string, error) {
	return vs.run("pitbull run -u " + passlistUrl + " -w " + walletString)
}

// RunPitbullForTokenlist - runs Pitbull process for given tokenlist and walletString.
func (vs *VastSSHClient) RunPitbullForTokenlist(tokenlist, walletString string) (string, error) {
	return vs.run("pitbull run -t '" + tokenlist + "' -w " + walletString)
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
			ssh.PublicKeys(vs.signer),
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
		output += line + "\n"
	}

	return output, nil
}

// AddKnownHost - dynamically registers host as known host.
// This method is partially working - in order to make it work properly, we would need to add
// dynamic keyscan for the host and handle known host entry generation.
func (vs *VastSSHClient) addKnownHost(sshDirPath, ipAddress string) (ssh.HostKeyCallback, error) {
	knownHostsPath := sshDirPath + "/known_hosts"
	publicKeyPath := sshDirPath + "/id_rsa.pub"

	publicKeyRaw, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	// The format that is used for storing SSH public key in a file is a different format
	// that we need for ParsePublicKey (wire format), therefore we need to parse it accordingly.
	publicKeyFileFormat, _, _, _, err := ssh.ParseAuthorizedKey(publicKeyRaw)
	if err != nil {
		return nil, err
	}

	publicKey, err := ssh.ParsePublicKey(publicKeyFileFormat.Marshal())
	if err != nil {
		return nil, err
	}

	hostNormalized := knownhosts.Normalize(ipAddress)

	knownHost := knownhosts.Line([]string{hostNormalized}, publicKey)

	knownHostFile, err := os.OpenFile(knownHostsPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	if _, err := knownHostFile.WriteString(knownHost + "\n"); err != nil {
		return nil, err
	}

	knownHostFile.Close()

	return knownhosts.New(knownHostsPath)
}
