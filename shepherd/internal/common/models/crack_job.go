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

	Keyword      string `bson:"keyword" json:"keyword"`
	WalletString string `bson:"walletString" json:"walletString"`
	PasslistUrl  string `bson:"passlistUrl" json:"passlistUrl"`

	InstanceId primitive.ObjectID `bson:"instanceId" json:"instanceId"`
	Instance   *PitbullInstance   `bson:"instance,omitempty" json:"instance"`

	StartedAt NullableTime `bson:"startedAt" json:"startedAt"`

	ErrorLog string `bson:"errorLog" json:"errorLog"`
}

// NewPitbullJob - returns new PitbullJob instance.
func NewPitbullJob(keyword, passlistUrl, walletString string) *CrackJob {
	job := &CrackJob{
		Keyword:      keyword,
		WalletString: walletString,
		PasslistUrl:  passlistUrl,
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

// PitbullJobsListPayload - describes a payload for PitbullJobs list.
type PitbullJobsListPayload struct {
	api.BaseListPayload

	Statuses []JobStatusEnum
	Keyword  string
}

// NewPitbullJobsListPayload - returns new PitbullJobsListPayload instance.
func NewPitbullJobsListPayload() *PitbullJobsListPayload {
	return &PitbullJobsListPayload{}
}

// FromRequest - populates payload values from request.
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

	pj.Statuses = statuses
	pj.Page = page
	pj.PageSize = pageSize
	pj.Keyword = keyword

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
