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

	Name         string   `bson:"name" json:"name"`
	Keyword      string   `bson:"keyword" json:"keyword"`
	WalletString string   `bson:"walletString" json:"walletString"`
	PasslistUrl  string   `bson:"passlistUrl" json:"passlistUrl"`
	Tokens       []string `bson:"tokens" json:"tokens"`

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

// GetTokenlist - builds tokenlist from Tokens, by joining them with newlines ('\n').
func (pj *CrackJob) GetTokenlist() string {
	if len(pj.Tokens) == 0 {
		return ""
	}

	return strings.Join(pj.Tokens, "\n")
}

// PitbullJobsListPayload - describes a payload for PitbullJobs list.
type PitbullJobsListPayload struct {
	api.BaseListPayload

	Statuses []JobStatusEnum
	JobId    string
	Keyword  string
	Name     string
}

// NewPitbullJobsListPayload - returns new PitbullJobsListPayload instance.
func NewPitbullJobsListPayload() *PitbullJobsListPayload {
	return &PitbullJobsListPayload{}
}

// Populate - populates payload values from request.
func (pj *PitbullJobsListPayload) Populate(r *http.Request) error {
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
	jobId := r.URL.Query().Get("jobId")
	name := r.URL.Query().Get("name")

	pj.Statuses = statuses
	pj.Page = page
	pj.PageSize = pageSize
	pj.Keyword = keyword
	pj.Name = name
	pj.JobId = jobId

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
