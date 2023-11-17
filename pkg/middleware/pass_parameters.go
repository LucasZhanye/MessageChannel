package middleware

import "github.com/gin-gonic/gin"

func PassParameters(key string, value any) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(key, value)
	}
}
