package middleware

import (
	"messagechannel/internal/pkg/core/server/response"
	"messagechannel/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWT(secret string) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		resp := response.OK

		token := ctx.Request.Header.Get("token")
		if token == "" {
			resp = response.ErrUnAuth.WithMessage("token is null")
		} else {
			claims, err := utils.ParseToken(token, []byte(secret))
			if err != nil {
				resp = response.ErrUnAuth.WithMessage(err.Error())
			} else {
				ctx.Set("username", claims.UserName)
			}
		}

		if resp.Code != 200 {
			ctx.JSON(http.StatusUnauthorized, resp)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
