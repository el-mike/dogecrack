package vast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/el-mike/dogecrack/shepherd/provider"
)

const (
	// Please note that we need to drop ' signs at the beginning and end -
	// exec.Command will add them for us.
	CheapOfferFilter = "dph < 0.3"
)

// VastCLIClient - facade for vast CLI operations.
type VastCLIClient struct {
	apiSecret    string
	pitbullImage string
}

// NewVastCli - returns new VastCli instance.
func NewVastCLI(apiSecret, pitbullImage string) *VastCLIClient {
	return &VastCLIClient{
		apiSecret:    apiSecret,
		pitbullImage: pitbullImage,
	}
}

// GetInstances - returns current instances.
func (vc *VastCLIClient) GetInstances() ([]*VastInstance, error) {
	result, err := vc.run("show", "instances")
	if err != nil {
		return nil, err
	}

	return vc.parseInstances(result)
}

// StartInstance - starts new Vast.ai instance. Waits for starting process to be over.
func (vc *VastCLIClient) StartInstance(offerId int) (*VastCreateResponse, error) {
	result, err := vc.run("create", "instance", strconv.Itoa(offerId), "--image", vc.pitbullImage)
	if err != nil {
		return nil, err
	}

	res := string(result)

	fmt.Print(res)

	var response *VastCreateResponse

	if err := json.Unmarshal(result, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetInstance - returns single, existing (rented) instance based on passed id.
func (vc *VastCLIClient) GetInstance(instanceId int) (*VastInstance, error) {
	instances, err := vc.GetInstances()
	if err != nil {
		return nil, err
	}

	for _, instance := range instances {
		if instance.ID == instanceId {
			return instance, nil
		}
	}

	return nil, provider.NewInstanceNotFound(instanceId)
}

// GetOfferByCriteria - returns first offer matching the criteria.
func (vc *VastCLIClient) GetOfferByCriteria(criteria string) (*VastOffer, error) {
	offers, err := vc.GetOffers(criteria)
	if err != nil {
		return nil, err
	}

	return offers[0], nil
}

// GetOffers - returns current Vast.ai machine offers.
func (vc *VastCLIClient) GetOffers(filter string) ([]*VastOffer, error) {
	result, err := vc.run("search", "offers", filter)
	if err != nil {
		return nil, err
	}

	var offers []*VastOffer

	if err := json.Unmarshal(result, &offers); err != nil {
		return nil, err
	}

	return offers, nil
}

// run - runs single command.
func (vc *VastCLIClient) run(cmdArgs ...string) ([]byte, error) {
	// We append apiKey and output format args.
	cmdArgs = append(cmdArgs, "--api-key", vc.apiSecret, "--raw")

	cmd := exec.Command("vast", cmdArgs...)

	var out bytes.Buffer
	cmd.Stdout = &out

	var errOut bytes.Buffer
	cmd.Stderr = &errOut

	if err := cmd.Run(); err != nil {
		fmt.Println(string(out.Bytes()))
		fmt.Print(string(errOut.Bytes()))
		return nil, err
	}

	return out.Bytes(), nil
}

// parseInstances - helper function for parsing cmd result as VastInstance slice.
func (vc *VastCLIClient) parseInstances(cmdResult []byte) ([]*VastInstance, error) {
	var instances []*VastInstance

	if err := json.Unmarshal(cmdResult, &instances); err != nil {
		return nil, err
	}

	return instances, nil
}
