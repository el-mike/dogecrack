package vast

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

// GetFakeVastIp - returns an ID of fake_vast container when running locally.
func GetFakeVastIp(rootDir string, fakeVastId int) (string, error) {
	cmd := exec.Command(rootDir+"/tools/vast/scripts/get_fake_vast_ip.sh", strconv.Itoa(fakeVastId))

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Sice echo command retirns a newline at the end, we want to
	// make sure ip is correctly trimmed.
	ip := strings.Trim(out.String(), "\n")

	return ip, nil
}

// AddSSHFingerprint - adds SSH fingerprint to host's known_hosts file, to prevent
// "host uknown" errors while connection to remote machine.
func AddSSHFingerprint(rootDir, host, sshDirPath string) error {
	hostsFilePath := sshDirPath + "/known_hosts"

	cmd := exec.Command(rootDir+"/tools/vast/scripts/add_ssh_fingerprint.sh", host, hostsFilePath)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
