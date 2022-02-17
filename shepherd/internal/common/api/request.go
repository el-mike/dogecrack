package api

import (
	"net/http"
	"strconv"
)

// BaseListPayload - describes basic parameters passed in list requests.
type BaseListPayload struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func GetBaseListPayload(r *http.Request) (int, int, error) {
	pageParam := r.URL.Query().Get("page")
	pageSizeParam := r.URL.Query().Get("pageSize")

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		return 0, 0, err
	}

	pageSize, err := strconv.Atoi(pageSizeParam)
	if err != nil {
		return 0, 0, err
	}

	return page, pageSize, nil
}
