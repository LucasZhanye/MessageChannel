package protocol

import "fmt"

var _ error = (*Error)(nil)

type Error struct {
	Code    int
	Message string
	Reason  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("code = %d, msg = %s, reason = %s", e.Code, e.Message, e.Reason)
}

func (e *Error) WithReason(reason string) *Error {
	e.Reason = reason
	return e
}

var (
	ErrConnect = &Error{
		Code:    10001,
		Message: "connect fail",
	}

	ErrClientIdMatch = &Error{
		Code:    10002,
		Message: "client identifie not match",
	}

	ErrTopicName = &Error{
		Code:    10003,
		Message: "topic not match rule",
	}

	ErrGroupIsNull = &Error{
		Code:    10004,
		Message: "group is null",
	}

	ErrSubscribe = &Error{
		Code:    10005,
		Message: "subscribe fail",
	}

	ErrUnsubscribe = &Error{
		Code:    10006,
		Message: "unsubscribe fail",
	}

	ErrDataIsNull = &Error{
		Code:    10007,
		Message: "data is null",
	}

	ErrPublish = &Error{
		Code:    10008,
		Message: "publish fail",
	}

	ErrAck = &Error{
		Code:    10009,
		Message: "publish fail",
	}

	ErrSyncPublishReply = &Error{
		Code:    10010,
		Message: "sync publish reply fail",
	}
)
