package model

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	UserName string `json:"username"`
	Token    string `json:"token"`
}

type SystemInfo struct {
	TotalClient       int    `json:"total_client"`
	TotalSubscription int    `json:"total_subscription"`
	Name              string `json:"name"`
	Version           string `json:"version"`
	Engine            string `json:"engine"`
}

type ClientSubscriptions struct {
	Topic string `json:"topic"`
	Group string `json:"group"`
}

type ClientBaseInfo struct {
	Identifie   string `json:"identifie"`
	Address     string `json:"address"`
	ConnectTime int64  `json:"connect_time"`
}

type ClientInfo struct {
	ClientBaseInfo
	SubscriptionInfo []*ClientSubscriptions `json:"subscription_info"`
}

type TopicInfo struct {
	TopicName string       `json:"topic_name"`
	Group     []*GroupInfo `json:"group"`
}

type GroupInfo struct {
	GroupName    string            `json:"group_name"`
	Subscription []*ClientBaseInfo `json:"subscription"`
}

type Subscription struct {
	Identifie string `json:"identifie" binding:"required"`
	Topic     string `json:"topic" binding:"required"`
	Group     string `json:"group" binding:"required"`
}


