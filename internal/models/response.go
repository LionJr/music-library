package models

type Response struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Message string `json:"msg"`
}

type ErrorResponseWithCode struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"msg"`
}
