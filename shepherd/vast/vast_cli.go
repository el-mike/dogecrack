package vast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

// VastCLI - facade for vast CLI operations.
type VastCLI struct {
	apiSecret string
}

// NewVastCli - returns new VastCli instance.
func NewVastCLI(apiSecret string) *VastCLI {
	return &VastCLI{
		apiSecret: apiSecret,
	}
}

// GetInstances - returns current instances.
func (vc *VastCLI) GetInstances() ([]*VastInstance, error) {
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

// run - runs single command.
func (vc *VastCLI) run(cmdArgs ...string) ([]byte, error) {
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
