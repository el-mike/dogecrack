package api

import "encoding/json"

// ListResponse - struct describing response object for lists.
type ListResponse struct {
	Data       json.RawMessage `json:"data"`
	TotalCount int             `json:"totalCount"`
	Page       int             `json:"page"`
}

// NewListResponse - returns new ListResponse object.
func NewListResponse(data json.RawMessage, page, totalCount int) *ListResponse {
	return &ListResponse{
		Data:       data,
		TotalCount: totalCount,
		Page:       page,
	}
}
