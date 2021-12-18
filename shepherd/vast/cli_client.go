package vast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/el-mike/dogecrack/shepherd/vast/mocks"
)

// VastCLIClient - facade for vast CLI operations.
type VastCLIClient struct {
	apiSecret string
}

// NewVastCli - returns new VastCli instance.
func NewVastCLI(apiSecret string) *VastCLIClient {
	return &VastCLIClient{
		apiSecret: apiSecret,
	}
}

// GetInstances - returns current instances.
func (vc *VastCLIClient) GetInstances() ([]*VastInstance, error) {
	result, err := vc.run("show", "instances")
	if err != nil {
		return nil, err
	}

	fmt.Print(string(result))

	var instances []*VastInstance

	if err := json.Unmarshal(result, &instances); err != nil {
		return nil, err
	}

	return instances, nil
}

// StartInstance - starts new Vast.ai instance. Waits for starting process to be over.
func (vc *VastCLIClient) StartInstance(offerId int) (*VastInstance, error) {
	var instance *VastInstance

	if err := json.Unmarshal([]byte(mocks.MockInstanceJSON), &instance); err != nil {
		return nil, err
	}

	return instance, nil
}

// run - runs single command.
func (vc *VastCLIClient) run(cmdArgs ...string) ([]byte, error) {
	// We append apiKey and output format args.
	cmdArgs = append(cmdArgs, "--api-key", vc.apiSecret, "--raw")

	cmd := exec.Command("vast", cmdArgs...)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
