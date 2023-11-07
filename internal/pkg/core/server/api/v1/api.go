package v1

import (
	"messagechannel/internal/pkg/core"
	"messagechannel/internal/pkg/version"
	"net/http"

	"github.com/gin-gonic/gin"
)

func VersionHandler(ctx *gin.Context) {
	ctx.JSON(200, version.Get())
}

func ClientInfoHandler(ctx *gin.Context) {
	node, exist := ctx.Get("node")
	if !exist {
		ctx.JSON(http.StatusInternalServerError, "Error")
		return
	}

	n := node.(*core.Node)
	clients := n.GetClientManager().GetAll()

	ctx.JSON(http.StatusOK, clients)
}

func SubscriptionHandler(ctx *gin.Context) {
	node, exist := ctx.Get("node")
	if !exist {
		ctx.JSON(http.StatusInternalServerError, "Error")
		return
	}

	n := node.(*core.Node)
	subs := n.GetSubscriptionManager().GetAll()

	ctx.JSON(http.StatusOK, subs)
}
