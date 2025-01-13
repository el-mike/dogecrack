package vast

// SearchCriteriaProvider - an entity responsible for providing vast.ai search criteria.
type SearchCriteriaProvider interface {
	// GetSearchCriteria - returns search criteria Vast.ai CLI should use to get instances.
	GetSearchCriteria() string
}
