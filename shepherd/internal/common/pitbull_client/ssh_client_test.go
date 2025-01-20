package pitbull_client

import (
	"github.com/el-mike/dogecrack/shepherd/internal/config"
	"github.com/el-mike/dogecrack/shepherd/internal/vast"
	"path/filepath"
	"strings"
	"testing"
)

// TestSSHConnectionFakeVast - tests SSH connection.
// Requires fake Vast.ai instance to be run.
func TestSSHConnectionFakeVast(t *testing.T) {
	// t.Skip()

	rootPath, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}

	appConfig, err := config.NewAppConfig(rootPath)
	if err != nil {
		panic(err)
	}

	fakeVastOneIpAddress, err := vast.GetFakeVastIp(rootPath, 1)
	if err != nil {
		panic(err)
	}

	fakeVastOnePort := 22

	client, err := NewPitbullSSHClient(appConfig.SSHUser, appConfig.SSHPassword, appConfig.SSHDirPath, appConfig.SSHPrivateKey, fakeVastOneIpAddress, fakeVastOnePort)
	if err != nil {
		panic(err)
	}

	result, err := client.GetPitbullStatus()
	if err != nil {
		panic(err)
	}

	if strings.TrimSpace(result) != "WAITING" {
		t.Fatalf("Incorrect status: %s", result)
	}
}
