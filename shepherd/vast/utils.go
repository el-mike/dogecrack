package vast

import (
	"bytes"
	"os/exec"
	"strings"
)

// GetFakeVastIp - returns an ID of fake_vast container when running locally.
func GetFakeVastIp(rootDir string) string {
	cmd := exec.Command(rootDir + "/vast/scripts/get_fake_vast_ip.sh")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	// Sice echo command retirns a newline at the end, we want to
	// make sure ip is correctly trimmed.
	ip := strings.Trim(out.String(), "\n")

	return ip
}
