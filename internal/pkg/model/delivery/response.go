package delivery

// This file is only used for swagger docs

// APIResponse
// @Description API response wrapper
type APIResponse struct {
	Status int         `json:"status" example:"200" description:"HTTP status code"`
	Body   interface{} `json:"body" description:"Response data"`
}

type APIErrorResponse struct {
	Status int    `json:"status" example:"400" description:"HTTP status code"`
	Error  string `json:"error" example:"Something went wrong" description:"Error message"`
}

// APIBadRequestErrorResponse
// @Description API bad request error response structure
type APIBadRequestErrorResponse struct {
	Status int    `json:"status" example:"400" description:"HTTP status code"`
	Error  string `json:"error" example:"Something went wrong" description:"Error message"`
}

// APIInternalServerErrorResponse
// @Description API internal server error response structure
type APIInternalServerErrorResponse struct {
	Status int    `json:"status" example:"500" description:"HTTP status code"`
	Error  string `json:"error" example:"Something went wrong" description:"Error message"`
}
