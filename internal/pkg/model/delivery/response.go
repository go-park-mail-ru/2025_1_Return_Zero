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

// Message
// @Description Message for responses without data
type Message struct {
	Message string `json:"msg" example:"object have been successfully created/updated" description:"Message for responses without data"`
}

// APIUnauthorizedErrorResponse
// @Description API unauthorized error response structure
type APIUnauthorizedErrorResponse struct {
	Status int    `json:"status" example:"401" description:"HTTP status code"`
	Error  string `json:"error" example:"Unauthorized" description:"Error message"`
}

// APIForbiddenErrorResponse
// @Description API forbidden error response structure
type APIForbiddenErrorResponse struct {
	Status int    `json:"status" example:"403" description:"HTTP status code"`
	Error  string `json:"error" example:"Forbidden" description:"Error message"`
}

// APIRequestEntityTooLargeErrorResponse
// @Description API request entity too large error response structure
type APIRequestEntityTooLargeErrorResponse struct {
	Status int    `json:"status" example:"413" description:"HTTP status code"`
	Error  string `json:"error" example:"Request entity too large" description:"Error message"`
}

// APINotFoundErrorResponse
// @Description API not found error response structure
type APINotFoundErrorResponse struct {
	Status int    `json:"status" example:"404" description:"HTTP status code"`
	Error  string `json:"error" example:"Not found" description:"Error message"`
}
