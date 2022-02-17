package models

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/el-mike/dogecrack/shepherd/internal/common/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JobStatus - describes the status of a single PitbullJob.
type JobStatus int8

const (
	JobScheduled JobStatus = iota
	JobProcessing
	JobRescheduled
	JobRejected
	JobAcknowledged
)

var nameByJobStatus = map[JobStatus]string{
	JobScheduled:    "JOB_SCHEDULED",
	JobProcessing:   "JOB_PROCESSING",
	JobRescheduled:  "JOB_RESCHEDULED",
	JobRejected:     "JOB_REJECTED",
	JobAcknowledged: "JOB_ACKNOWLEDGED",
}

// PitbullJob - represents a single Pitbull job.
type PitbullJob struct {
	BaseModel `bson:",inline"`

	Keyword      string `bson:"keyword" json:"keyword"`
	WalletString string `bson:"walletString" json:"walletString"`
	PasslistUrl  string `bson:"passlistUrl" json:"passlistUrl"`

	InstanceId primitive.ObjectID `bson:"instanceId" json:"instanceId"`
	Instance   *PitbullInstance   `bson:"instance,omitempty" json:"instance"`

	StartedAt NullableTime `bson:"startedAt" json:"startedAt"`

	FirstScheduledAt NullableTime `bson:"firstScheduledAt" json:"firstScheduledAt,omitempty"`
	LastScheduledAt  NullableTime `bson:"lastScheduledAt" json:"lastScheduledAt,omitempty"`
	AcknowledgedAt   NullableTime `bson:"acknowledgedAt" json:"acknowledgedAt,omitempty"`
	RejectedAt       NullableTime `bson:"rejectedAt" json:"rejectedAt,omitempty"`

	Status          JobStatus `bson:"status" json:"status"`
	RescheduleCount int       `bson:"rescheduleCount" json:"rescheduleCount"`

	ErrorLog string `bson:"errorLog" json:"errorLog"`
}

// NewPitbullJob - returns new PitbullJob instance.
func NewPitbullJob(keyword, passlistUrl, walletString string) *PitbullJob {
	job := &PitbullJob{
		Keyword:      keyword,
		WalletString: walletString,
		PasslistUrl:  passlistUrl,
		Status:       JobScheduled,
	}

	job.ID = primitive.NewObjectID()

	return job
}

// AppendError - appends an error to PitbullJob's internal error log.
func (pj *PitbullJob) AppendError(err error) {
	if pj.ErrorLog == "" {
		pj.ErrorLog = ""
	}

	pj.ErrorLog += fmt.Sprintf("%s\n", err.Error())
}

// PitbullJobsListPayload - describes a payload for PitbullJobs list.
type PitbullJobsListPayload struct {
	api.BaseListPayload

	Statuses []JobStatus
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
	statuses := []JobStatus{}

	if statusesParam != "" {
		statusesRaw := strings.Split(statusesParam, ",")

		for _, statusRaw := range statusesRaw {
			status, err := strconv.Atoi(statusRaw)
			if err != nil {
				return err
			}

			statuses = append(statuses, JobStatus(status))
		}
	}

	keyword := r.URL.Query().Get("keyword")

	pj.Statuses = statuses
	pj.Page = page
	pj.PageSize = pageSize
	pj.Keyword = keyword

	return nil
}

// PagedPitbullJobs - paged DB result from PitbullJobs.
type PagedPitbullJobs struct {
	Data     []*PitbullJob `bson:"data" json:"data"`
	PageInfo *PageInfo     `bson:"pageInfo" json:"pageInfo"`
}

// NewPagedPitbullJobs - returns new PagedPitbullJobs instance.
func NewPagedPitbullJobs() *PagedPitbullJobs {
	return &PagedPitbullJobs{
		Data: []*PitbullJob{},
		PageInfo: &PageInfo{
			Page:     0,
			PageSize: 0,
			Total:    0,
		},
	}
}
