package dto

type SuccessResp200 struct {
	Success    bool `json:"success" default:"true"`
	StatusCode int  `json:"statusCode" default:"200"`
	Data       any  `json:"data"`
}

type SuccessResp201 struct {
	Success    bool `json:"success" default:"true"`
	StatusCode int  `json:"statusCode" default:"201"`
	Data       any  `json:"data"`
}

type ErrorMsgResp struct {
	Message string `json:"message"`
}

type FailedResp400 struct {
	Success    bool           `json:"success" default:"false"`
	StatusCode int            `json:"statusCode" default:"400"`
	Errors     []ErrorMsgResp `json:"errors"`
}

type FailedResp401 struct {
	Success    bool           `json:"success" default:"false"`
	StatusCode int            `json:"statusCode" default:"401"`
	Errors     []ErrorMsgResp `json:"errors"`
}

type FailedResp404 struct {
	Success    bool           `json:"success" default:"false"`
	StatusCode int            `json:"statusCode" default:"404"`
	Errors     []ErrorMsgResp `json:"errors"`
}

type FailedResp422 struct {
	Success    bool           `json:"success" default:"false"`
	StatusCode int            `json:"statusCode" default:"422"`
	Errors     []ErrorMsgResp `json:"errors"`
}

type FailedResp500 struct {
	Success    bool           `json:"success" default:"false"`
	StatusCode int            `json:"statusCode" default:"500"`
	Errors     []ErrorMsgResp `json:"errors"`
}

type DocPaginationResp struct {
	Total      int  `json:"total"`
	IsLoadMore bool `json:"isLoadMore"`
	Data       any  `json:"data"`
	Next       Next `json:"next"`
	Prev       Prev `json:"prev"`
}
