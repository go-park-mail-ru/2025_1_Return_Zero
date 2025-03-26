package model

// Pagination represents pagination parameters for API requests
// @Description Pagination parameters for list endpoints
type Pagination struct {
	Offset int `json:"offset" example:"0" description:"Number of items to skip"`
	Limit  int `json:"limit" example:"10" description:"Maximum number of items to return"`
}

const PaginationKey = "pagination"
