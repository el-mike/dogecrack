package pitbull_client

import (
	"bufio"
	"fmt"
	"github.com/el-mike/dogecrack/shepherd/internal/common"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

const CONN_PROTOCOL = "tcp"

// PitbullSSHClient - SSH connection client for running Pitbull.
type PitbullSSHClient struct {
	user      string
	password  string
	signer    ssh.Signer
	hostKeyCb ssh.HostKeyCallback

	ipAddress string
	port      int

	conn *ssh.Client

	logger *common.Logger
}

// NewPitbullSSHClient - returns new PitbullSSHClient instance. Does NOT start up a connection.
func NewPitbullSSHClient(user, password, sshDirPath, sshPrivateKey, ipAddress string, port int) (*PitbullSSHClient, error) {
	// We store sshPrivateKey as string with '\n', but since it's coming from .env file, it treats
	// '\n' characters as normal signs. Therefore, we use replace to insert actual new lines.
	sshPrivateKey = strings.ReplaceAll(sshPrivateKey, `\n`, "\n")
	privateKeyRaw := []byte(sshPrivateKey)

	signer, err := ssh.ParsePrivateKey(privateKeyRaw)
	if err != nil {
		return nil, err
	}

	client := &PitbullSSHClient{
		user:     user,
		password: password,
		signer:   signer,
		// TODO: Consider handling host key verification - it's partially implemented below.
		hostKeyCb: ssh.InsecureIgnoreHostKey(),
		ipAddress: ipAddress,
		port:      port,

		logger: common.NewLogger("SSHClient", os.Stdout, os.Stderr),
	}

	return client, nil
}

// RunPitbullForPasslist - runs Pitbull process for given passlistUrl and walletString.
func (vs *PitbullSSHClient) RunPitbullForPasslist(walletString, passlistUrl string, skipCount, minLength, maxLength int64) (string, error) {
	cmd := BuildRunCommand(walletString, passlistUrl, "", skipCount, minLength, maxLength)

	vs.logger.Info.Printf("Running Pitbull command for passlist: \"%s\"", cmd)
	return vs.Run(cmd)
}

// RunPitbullForTokenlist - runs Pitbull process for given tokenlist and walletString.
func (vs *PitbullSSHClient) RunPitbullForTokenlist(walletString, tokenlist string, skipCount, minLength, maxLength int64) (string, error) {
	cmd := BuildRunCommand(walletString, "", tokenlist, skipCount, minLength, maxLength)

	vs.logger.Info.Printf("Running Pitbull command for tokenlist: \"%s\"", cmd)
	return vs.Run(cmd)
}

// GetPitbullStatus - runs Pitbull's status command and returns the output.
func (vs *PitbullSSHClient) GetPitbullStatus() (string, error) {
	return vs.Run("pitbull status")
}

// GetPitbullStatus - runs Pitbull's progress command and returns the output.
func (vs *PitbullSSHClient) GetPitbullProgress() (string, error) {
	return vs.Run("pitbull progress")
}

// GetPitbullOutput - runs Pitbull's output command and returns the output.
func (vs *PitbullSSHClient) GetPitbullOutput() (string, error) {
	return vs.Run("pitbull output")
}

// Connect - starts a SSH connection.
func (vs *PitbullSSHClient) connect() error {
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
func (vs *PitbullSSHClient) close() error {
	return vs.conn.Close()
}

func (vs *PitbullSSHClient) Run(cmd string) (string, error) {
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
func (vs *PitbullSSHClient) addKnownHost(sshDirPath, ipAddress string) (ssh.HostKeyCallback, error) {
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
