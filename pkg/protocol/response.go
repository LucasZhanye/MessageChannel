package protocol

var (
	OK        = response("success")
	ErrParam  = response("param error")
	ErrServer = response("internal error")
)

type ResponseType uint8

const (
	RESPONSE_CONNECT ResponseType = iota + 1
	RESPONSE_PUBLISH
	RESPONSE_SUBSCRIBE
	RESPONSE_ACK
)

// Response sent from server to client
type Response struct {
	Type     ResponseType `json:"type"`
	Reason   string       `json:"reason"`
	Message  string       `json:"message"`
	MetaData any          `json:"metadata"`
}

func response(reason string) *Response {
	return &Response{
		Reason: reason,
	}
}

func (resp *Response) WithType(rtype ResponseType) *Response {
	resp.Type = rtype
	return resp
}

func (resp *Response) WithMessage(message string) *Response {
	resp.Message = message
	return resp
}

func (resp *Response) WithMetaData(data any) *Response {
	resp.MetaData = data
	return resp
}
