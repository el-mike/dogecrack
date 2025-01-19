package models

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/el-mike/dogecrack/shepherd/internal/common/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CrackJob - represents a single Pitbull job.
type CrackJob struct {
	Job `bson:",inline"`

	Name                  string                    `bson:"name" json:"name"`
	Keyword               string                    `bson:"keyword" json:"keyword"`
	WalletString          string                    `bson:"walletString" json:"walletString"`
	PasslistUrl           string                    `bson:"passlistUrl" json:"passlistUrl"`
	Tokenlist             string                    `bson:"tokenlist" json:"tokenlist"`
	TokenGeneratorVersion TokenGeneratorVersionEnum `bson:"tokenGeneratorVersion" json:"tokenGeneratorVersion"`

	InstanceId primitive.ObjectID `bson:"instanceId" json:"instanceId"`
	Instance   *PitbullInstance   `bson:"instance,omitempty" json:"instance"`

	StartedAt NullableTime `bson:"startedAt" json:"startedAt"`

	ErrorLog string `bson:"errorLog" json:"errorLog"`
}

// NewCrackJob - returns new PitbullJob instance.
func NewCrackJob(walletString string) *CrackJob {
	job := &CrackJob{
		WalletString: walletString,
	}

	job.ID = primitive.NewObjectID()
	job.Status = JobStatus.Scheduled

	return job
}

// AppendError - appends an error to PitbullJob's internal error log.
func (pj *CrackJob) AppendError(err error) {
	if pj.ErrorLog == "" {
		pj.ErrorLog = ""
	}

	pj.ErrorLog += fmt.Sprintf("%s\n", err.Error())
}

// CrackJobsListPayload - describes a payload for PitbullJobs list.
type CrackJobsListPayload struct {
	api.BaseListPayload

	Statuses              []JobStatusEnum
	JobId                 string
	Keyword               string
	PasslistUrl           string
	Name                  string
	TokenGeneratorVersion TokenGeneratorVersionEnum
}

// NewCrackJobsListPayload - returns new CrackJobsListPayload instance.
func NewCrackJobsListPayload() *CrackJobsListPayload {
	return &CrackJobsListPayload{}
}

// Populate - populates payload values from request.
func (pj *CrackJobsListPayload) Populate(r *http.Request) error {
	page, pageSize, err := api.GetBaseListPayload(r)
	if err != nil {
		return err
	}

	statusesParam := strings.Trim(r.URL.Query().Get("statuses"), ",")
	statuses := []JobStatusEnum{}

	if statusesParam != "" {
		statusesRaw := strings.Split(statusesParam, ",")

		for _, statusRaw := range statusesRaw {
			status, err := strconv.Atoi(statusRaw)
			if err != nil {
				return err
			}

			statuses = append(statuses, JobStatusEnum(status))
		}
	}

	keyword := r.URL.Query().Get("keyword")
	passlistUrl := r.URL.Query().Get("passlistUrl")
	jobId := r.URL.Query().Get("jobId")
	name := r.URL.Query().Get("name")
	tokenGeneratorVersion := r.URL.Query().Get("tokenGeneratorVersion")

	pj.Statuses = statuses
	pj.Page = page
	pj.PageSize = pageSize
	pj.Keyword = keyword
	pj.PasslistUrl = passlistUrl
	pj.Name = name
	pj.JobId = jobId

	if tokenGeneratorVersion != "" {
		tokenGeneratorVersionInt, err := strconv.Atoi(tokenGeneratorVersion)
		if err != nil {
			return err
		}

		pj.TokenGeneratorVersion = TokenGeneratorVersionEnum(tokenGeneratorVersionInt)
	}

	return nil
}

// PagedCrackJobs - paged DB result from PitbullJobs.
type PagedCrackJobs struct {
	Data     []*CrackJob `bson:"data" json:"data"`
	PageInfo *PageInfo   `bson:"pageInfo" json:"pageInfo"`
}

// NewPagedCrackJobs - returns new PagedPitbullJobs instance.
func NewPagedCrackJobs() *PagedCrackJobs {
	return &PagedCrackJobs{
		Data: []*CrackJob{},
		PageInfo: &PageInfo{
			Page:     0,
			PageSize: 0,
			Total:    0,
		},
	}
}

type BaseCreateJobPayload struct {
	WalletString string
	Name         string
}

type CreateJobsForKeywordsPayload struct {
	BaseCreateJobPayload

	Keywords              []string
	Tokenlist             string
	TokenGeneratorVersion TokenGeneratorVersionEnum
}

type CancelCrackJobPayload struct {
	JobId string `json:"jobId"`
}

type RecreateCrackJobPayload struct {
	JobId string `json:"jobId"`
}
