package delivery

// Pagination represents pagination parameters for API requests
// @Description Pagination parameters for list endpoints
type Pagination struct {
	Offset int `json:"offset" mapstructure:"offset" example:"0" description:"Number of items to skip"`
	Limit  int `json:"limit" mapstructure:"limit" example:"10" description:"Maximum number of items to return"`
}

type PaginationConfig struct {
	MaxOffset     int `mapstructure:"max_offset"`
	MaxLimit      int `mapstructure:"max_limit"`
	DefaultOffset int `mapstructure:"default_offset"`
	DefaultLimit  int `mapstructure:"default_limit"`
}
