package vast

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/el-mike/dogecrack/shepherd/internal/common/host"
	"github.com/el-mike/dogecrack/shepherd/internal/vast/models"
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

	// Since vast.ai API is rather limited in terms of handling request,
	// we want to use mutual exclusion to perform request sequentially.
	sync.Mutex
}

// NewVastCli - returns new VastCli instance.
func NewVastCLI(apiSecret, pitbullImage string) *VastCLIClient {
	return &VastCLIClient{
		apiSecret:    apiSecret,
		pitbullImage: pitbullImage,
	}
}

// GetInstances - returns current instances.
func (vc *VastCLIClient) GetInstances() ([]*models.Instance, error) {
	vc.Lock()
	defer vc.Unlock()

	result, err := vc.run("show", "instances")
	if err != nil {
		return nil, err
	}

	return vc.parseInstances(result)
}

// StartInstance - starts new Vast.ai instance. Waits for starting process to be over.
func (vc *VastCLIClient) StartInstance(offerId int) (*models.CreateResponse, error) {
	vc.Lock()
	defer vc.Unlock()

	result, err := vc.run("create", "instance", strconv.Itoa(offerId), "--image", vc.pitbullImage)
	if err != nil {
		return nil, err
	}

	var response *models.CreateResponse

	if err := json.Unmarshal(result, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// DestroyInstance - stops a Vast.ai instance with given id.
func (vc *VastCLIClient) DestroyInstance(instanceId int) error {
	vc.Lock()
	defer vc.Unlock()

	_, err := vc.run("destroy", "instance", strconv.Itoa(instanceId))
	if err != nil {
		return err
	}

	return nil
}

// GetInstance - returns single, existing (rented) instance based on passed id.
func (vc *VastCLIClient) GetInstance(instanceId int) (*models.Instance, error) {
	instances, err := vc.GetInstances()
	if err != nil {
		return nil, err
	}

	for _, instance := range instances {
		if instance.ID == instanceId {
			return instance, nil
		}
	}

	return nil, host.NewHostInstanceNotFound(instanceId)
}

// GetOfferByCriteria - returns first offer matching the criteria.
func (vc *VastCLIClient) GetOfferByCriteria(criteria string) (*models.Offer, error) {
	offers, err := vc.GetOffers(criteria)
	if err != nil {
		return nil, err
	}

	return offers[0], nil
}

// GetOffers - returns current Vast.ai machine offers.
func (vc *VastCLIClient) GetOffers(filter string) ([]*models.Offer, error) {
	vc.Lock()
	defer vc.Unlock()

	result, err := vc.run("search", "offers", filter)
	if err != nil {
		return nil, err
	}

	var offers []*models.Offer

	if err := json.Unmarshal(result, &offers); err != nil {
		return nil, err
	}

	// If there are no offers, we want to return appropriate error, so the client
	// can act accordingly (for example with stopping jobs for some time).
	if len(offers) == 0 {
		return nil, host.NewNoHostMachinesAvailable()
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
		return nil, NewVastCLIError(string(out.Bytes()), string(errOut.Bytes()))
	}

	return vc.sanitizeJSONOutput(out.Bytes()), nil
}

// parseInstances - helper function for parsing cmd result as VastInstance slice.
func (vc *VastCLIClient) parseInstances(cmdResult []byte) ([]*models.Instance, error) {
	var instances []*models.Instance

	if err := json.Unmarshal(cmdResult, &instances); err != nil {
		return nil, err
	}

	return instances, nil
}

// sanitizeJSONOutput - makes sure vast CLI output has correct JSON format.
// This helps with any vast CLI additional messages that can be added to the output,
// as missing vast_pdf for example.
func (vc *VastCLIClient) sanitizeJSONOutput(out []byte) []byte {
	outString := string(out)
	jsonArrayStartIndex := strings.Index(outString, "[")
	jsonObjectStartIndex := strings.Index(outString, "{")

	jsonStartIndex := -1

	// Vast CLI can return both objects and arrays as root, depends on the command.
	// Therefore, we need to check for both cases, and prioritize accordingly.
	if jsonArrayStartIndex != -1 {
		jsonStartIndex = jsonArrayStartIndex
	}

	// If '{' was found, and it's before '[' (or '[' was not found), use it instead.
	if jsonObjectStartIndex != -1 && (jsonObjectStartIndex < jsonStartIndex || jsonStartIndex == -1) {
		jsonStartIndex = jsonObjectStartIndex
	}

	if jsonStartIndex == -1 {
		return out
	}

	jsonString := outString[jsonStartIndex:]

	return []byte(jsonString)
}
