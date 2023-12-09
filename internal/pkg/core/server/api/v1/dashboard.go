package v1

import (
	"messagechannel/internal/pkg/core"
	"messagechannel/internal/pkg/core/server/config"
	"messagechannel/internal/pkg/core/server/model"
	"messagechannel/internal/pkg/core/server/response"
	"messagechannel/internal/pkg/version"
	"messagechannel/pkg/protocol"
	"messagechannel/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const AdminName = "admin"
const GuestName = "guest"
const GuestPassword = "guest"

func LoginHandler(ctx *gin.Context) {
	user := &model.User{}

	node, exist := ctx.Get("node")
	if !exist {
		ctx.JSON(http.StatusOK, response.ErrServer)
		return
	}

	n := node.(*core.Node)

	cfg, exist := ctx.Get("config")
	if !exist {
		ctx.JSON(http.StatusOK, response.ErrServer)
		return
	}

	config := cfg.(*config.HttpConfig)

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusOK, response.ErrParam.WithMessage(err.Error()))
		return
	}

	if user.UserName == "" || user.Password == "" {
		ctx.JSON(http.StatusOK, response.ErrParam.WithMessage("UserName or Password is null"))
		return
	}

	// check user
	if user.UserName != AdminName && user.UserName != GuestName {
		resp := response.ErrParam.WithMessage("User not supported")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// check password
	if user.UserName == GuestName {
		ret := utils.SecureCompareString(user.Password, GuestPassword)
		if !ret {
			ctx.JSON(http.StatusOK, response.ErrParam.WithMessage("Password incorrect"))
			return
		}
	} else {
		// The inversion is done here, the user input is encry text, and the configuration file is plain text.
		// ret := utils.ComparePassword(user.Password, config.Web.Password)
		ret := utils.SecureCompareString(user.Password, config.Web.Password)
		if !ret {
			ctx.JSON(http.StatusOK, response.ErrParam.WithMessage("Password incorrect"))
			return
		}
	}
	// generate token
	token, err := utils.GenerateToken(&utils.Setting{
		Secret: []byte(config.Web.Secret),
		Customclaims: utils.Customclaims{
			UserName: user.UserName,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "messagechannel",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Web.ExpireTime)),
			},
		},
	})
	if err != nil {
		n.Log.Error("Generate token error: %v", err)
		ctx.JSON(http.StatusOK, response.ErrServer.WithMessage("Generate token fail"))
		return
	}

	userInfo := &model.UserInfo{
		UserName: user.UserName,
		Token:    token,
	}
	ctx.JSON(http.StatusOK, response.OK.WithMetaData(userInfo))
}

func SystemInfoHandler(ctx *gin.Context) {
	node, exist := ctx.Get("node")
	if !exist {
		ctx.JSON(http.StatusOK, response.ErrServer)
		return
	}

	n := node.(*core.Node)

	totalClient := n.GetClientManager().GetTotalClientCount()
	totalSubscription := n.GetSubscriptionManager().GetTotalSubscriptionCount()

	systemInfo := &model.SystemInfo{
		TotalClient:       totalClient,
		TotalSubscription: totalSubscription,
		Name:              n.Config.Name,
		Version:           version.Get(),
		Engine:            n.Config.EngineType,
	}

	ctx.JSON(http.StatusOK, response.OK.WithMetaData(systemInfo))
}

func ClientInfoHandler(ctx *gin.Context) {
	node, exist := ctx.Get("node")
	if !exist {
		ctx.JSON(http.StatusInternalServerError, "Error")
		return
	}

	n := node.(*core.Node)
	clients := n.GetClientManager().GetAll()

	clientInfos := make([]*model.ClientInfo, len(clients))

	for index, client := range clients {
		clientInfo := &model.ClientInfo{
			ClientBaseInfo: model.ClientBaseInfo{
				Identifie:   client.Identifie,
				Address:     client.Info.Address,
				ConnectTime: client.Info.ConnectTime,
			},
			SubscriptionInfo: []*model.ClientSubscriptions{},
		}

		if client.Info.SubInfo != nil {
			for _, sub := range client.Info.SubInfo.Items() {
				subscription := &model.ClientSubscriptions{
					Topic: sub.Topic,
					Group: sub.Group,
				}

				clientInfo.SubscriptionInfo = append(clientInfo.SubscriptionInfo, subscription)
			}
		}

		clientInfos[index] = clientInfo
	}

	ctx.JSON(http.StatusOK, response.OK.WithMetaData(clientInfos))
}

func SubscriptionHandler(ctx *gin.Context) {
	node, exist := ctx.Get("node")
	if !exist {
		ctx.JSON(http.StatusInternalServerError, "Error")
		return
	}

	n := node.(*core.Node)
	topics := n.GetSubscriptionManager().GetAll()

	topicInfo := make([]*model.TopicInfo, 0)

	topics.IterCb(func(topicName string, groups *core.Group) {
		tInfo := &model.TopicInfo{
			TopicName: topicName,
			Group:     make([]*model.GroupInfo, 0),
		}

		groups.IterCb(func(groupName string, subs *core.Subscription) {
			gInfo := &model.GroupInfo{
				GroupName:    groupName,
				Subscription: make([]*model.ClientBaseInfo, 0),
			}

			subs.IterCb(func(identifie string, clientInfo *core.Client) {
				subInfo := &model.ClientBaseInfo{
					Identifie:   clientInfo.Identifie,
					Address:     clientInfo.Info.Address,
					ConnectTime: clientInfo.Info.ConnectTime,
				}

				gInfo.Subscription = append(gInfo.Subscription, subInfo)
			})

			tInfo.Group = append(tInfo.Group, gInfo)
		})

		topicInfo = append(topicInfo, tInfo)
	})

	ctx.JSON(http.StatusOK, response.OK.WithMetaData(topicInfo))
}

func SubscribeHandler(ctx *gin.Context) {

	userName := ctx.GetString("username")
	if userName == AdminName {
		node, exist := ctx.Get("node")
		if !exist {
			ctx.JSON(http.StatusOK, response.ErrServer)
			return
		}

		n := node.(*core.Node)
		subscribeInfo := &model.Subscription{}
		if err := ctx.ShouldBindJSON(&subscribeInfo); err != nil {
			ctx.JSON(http.StatusOK, response.ErrParam.WithMessage(err.Error()))
			return
		}

		sub := protocol.NewSubscription(subscribeInfo.Identifie, subscribeInfo.Topic, subscribeInfo.Group)
		err := n.Subscribe(sub)
		if err != nil {
			ctx.JSON(http.StatusOK, response.ErrParam.WithMessage(err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, response.OK)
	} else {
		ctx.JSON(http.StatusOK, response.ErrPermission)
	}
}

func UnsubscribeHandler(ctx *gin.Context) {
	userName := ctx.GetString("username")
	if userName == AdminName {
		node, exist := ctx.Get("node")
		if !exist {
			ctx.JSON(http.StatusOK, response.ErrServer)
			return
		}

		n := node.(*core.Node)

		unsubscribeInfo := []*model.Subscription{}
		if err := ctx.ShouldBindJSON(&unsubscribeInfo); err != nil {
			ctx.JSON(http.StatusOK, response.ErrParam.WithMessage(err.Error()))
			return
		}
		failData := make([]*model.Subscription, 0)

		for _, unsub := range unsubscribeInfo {
			u := protocol.NewUnSubscription(unsub.Identifie, unsub.Topic, unsub.Group)
			err := n.UnSubscribe(u)
			if err != nil {
				failData = append(failData, unsub)
			}
		}

		if len(unsubscribeInfo) == len(failData) {
			ctx.JSON(http.StatusOK, response.ErrServer.WithMessage("all unsubscribe fail"))
			return
		}

		if len(failData) > 0 {
			ctx.JSON(http.StatusOK, response.OK.WithMetaData(failData))
			return
		}

		ctx.JSON(http.StatusOK, response.OK)
	} else {
		ctx.JSON(http.StatusOK, response.ErrPermission)
	}
}
