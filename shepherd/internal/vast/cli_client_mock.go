package vast

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/vast/models"
)

type fakeInstance struct {
	id             int
	ipAddress      string
	preStartChecks int
}

var fakeInstancesLimit = 2
var preStartChecksLimit = 3

var fakeVastOneId = 1111
var fakeVastTwoId = 2222

type instancesState []*fakeInstance

type VastCLIClientMock struct {
	rootDir string

	state instancesState
}

func NewVastCLIClientMock(rootDir string) *VastCLIClientMock {
	return &VastCLIClientMock{
		rootDir: rootDir,
		state:   instancesState{},
	}
}

// GetOfferByCriteria - returns first offer matching the criteria.
func (vc *VastCLIClientMock) GetOfferByCriteria(criteria string) (*models.Offer, error) {
	vc.simulateRequest()

	return &models.Offer{
		ID: 1111,
	}, nil
}

// GetInstance - returns single, existing (rented) instance based on passed id.
func (vc *VastCLIClientMock) GetInstance(instanceId int) (*models.Instance, error) {
	vc.simulateRequest()

	var instance *models.Instance

	for _, fakeInstance := range vc.state {
		if fakeInstance.id == instanceId {

			if fakeInstance.preStartChecks < preStartChecksLimit {
				fakeInstance.preStartChecks += 1
			}

			instance = vc.buildVastInstance(fakeInstance)
			break
		}
	}

	if instance == nil {
		return nil, errors.New(fmt.Sprintf("Instance: %d could not be found", instanceId))
	}

	return instance, nil
}

// StartInstance - starts new Vast.ai instance. Waits for starting process to be over.
func (vc *VastCLIClientMock) StartInstance(offerId int) (*models.CreateResponse, error) {
	vc.simulateRequest()

	if len(vc.state) >= fakeInstancesLimit {
		return nil, errors.New(fmt.Sprintf("There are only %d fake instances available", fakeInstancesLimit))
	}

	var id int
	var ipAddress string
	var err error

	if len(vc.state) == 0 {
		id = fakeVastOneId

		ipAddress, err = GetFakeVastIp(vc.rootDir, 1)
		if err != nil {
			return nil, err
		}

	} else if len(vc.state) == 1 {
		id = fakeVastTwoId

		ipAddress, err = GetFakeVastIp(vc.rootDir, 2)
		if err != nil {
			return nil, err
		}
	}

	vc.state = append(vc.state, &fakeInstance{
		id:             id,
		ipAddress:      ipAddress,
		preStartChecks: 0,
	})

	return &models.CreateResponse{
		Success:    true,
		InstanceId: id,
	}, nil
}

func (vc *VastCLIClientMock) DestroyInstance(instanceId int) error {
	if len(vc.state) == 0 {
		return nil
	}

	result := instancesState{}

	for _, instance := range vc.state {
		if instance.id != instanceId {
			result = append(result, instance)
		}
	}

	return nil
}

// GetInstances - returns current instances.
func (vc *VastCLIClientMock) GetInstances() ([]*models.Instance, error) {
	vc.simulateRequest()

	instances := []*models.Instance{}

	for _, fakeInstance := range vc.state {
		instances = append(instances, vc.buildVastInstance(fakeInstance))
	}

	return instances, nil
}

func (vc *VastCLIClientMock) simulateRequest() {
	timeout := rand.Intn(3) * int(time.Second)
	time.Sleep(time.Duration(timeout))
}

func (vc *VastCLIClientMock) buildVastInstance(fakeInstance *fakeInstance) *models.Instance {
	var status string

	if fakeInstance.preStartChecks < preStartChecksLimit {
		status = "loading"
	} else {
		status = "running"
	}

	return &models.Instance{
		ID:          fakeInstance.id,
		SSHHost:     fakeInstance.ipAddress,
		SSHPort:     22,
		Status:      status,
		DockerImage: "michalhuras/pitbull:latest",
	}
}
