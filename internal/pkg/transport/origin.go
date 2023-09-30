package transport

import (
	"net/http"
)

// CheckOriginFunc
func CheckOriginFunc() func(r *http.Request) bool {
	return func(r *http.Request) bool {
		return true
	}
}
