package models

type CreateResponse struct {
	Success    bool `json:"success"`
	InstanceId int  `json:"new_contract"`
}
