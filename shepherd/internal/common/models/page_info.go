package models

// PageInfo - describes a page of list data.
type PageInfo struct {
	Page     int `bson:"page" json:"page"`
	PageSize int `bson:"pageSize" json:"pageSize"`
	Total    int `bson:"total" json:"total"`
}
