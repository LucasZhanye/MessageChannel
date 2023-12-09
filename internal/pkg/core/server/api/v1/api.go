package v1

import (
	"messagechannel/internal/pkg/version"

	"github.com/gin-gonic/gin"
)

func VersionHandler(ctx *gin.Context) {
	ctx.JSON(200, version.Get())
}
