package transport

// CloseEvent
// Code reference https://developer.mozilla.org/zh-CN/docs/Web/API/CloseEvent
// 3000-3999 use by library or framework
// 4000-4999 use by application
type CloseEvent struct {
	Code   int
	Reason string
}

var CloseEventNormalClose = CloseEvent{
	Code:   3000,
	Reason: "normal close connect",
}

var CloseEventWriteError = CloseEvent{
	Code:   3001,
	Reason: "write error",
}
