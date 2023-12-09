package response

import "net/http"

type Response struct {
	Code     int         `json:"code"`
	Reason   string      `json:"reason"`
	Message  string      `json:"message"`
	MetaData interface{} `json:"metadata"`
}

func response(code int, reason string) *Response {
	return &Response{
		Code:   code,
		Reason: reason,
	}
}

func (resp *Response) WithMessage(message string) *Response {
	respCopy := *resp
	respCopy.Message = message
	return &respCopy
}

func (resp *Response) WithMetaData(data interface{}) *Response {
	respCopy := *resp
	respCopy.MetaData = data
	return &respCopy
}

var (
	OK            = response(http.StatusOK, "Success")
	ErrParam      = response(http.StatusBadRequest, "Params Error")
	ErrServer     = response(http.StatusInternalServerError, "Internal Server Error")
	ErrUnAuth     = response(http.StatusUnauthorized, "User unauthorized")
	ErrPermission = response(http.StatusForbidden, "Permission denied")
)
