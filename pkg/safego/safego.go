package safego

import (
	"fmt"
	"messagechannel/pkg/logger"
	"runtime"
)

func Execute(log logger.Log, fn func()) {
	go func(log logger.Log) {
		defer recoverPanic(log)
		fn()
	}(log)
}

func recoverPanic(log logger.Log) {
	if r := recover(); r != nil {
		buf := make([]byte, 2048)
		n := runtime.Stack(buf, false)
		stackInfo := fmt.Sprintf("%s", buf[:n])
		buf = nil

		log.Error("panic stack info %s, error: %v", stackInfo, r)
	}
}
