package config

import (
	"crypto/tls"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

const moduleName = "engine.rabbitmq"

const (
	BASE_TOPIC           = "default.basic.topic"
	DEAD_LETTER_EXCHANGE = "default.deadletter.exchange"
	DEAD_LETTER_QUEUE    = "default.deadletter.queue"
	DEAD_LETTER_KEY      = "default.deadletter.key"
	RESPONSE_EXCHANGE    = "default.response.exchange"
)

type RabbitMqConfig struct {
	Url        string
	AmqpConfig *amqp.Config
	TlsConfig  *tls.Config
	Exchange   *ExchangeConfig
	Queue      *QueueConfig
	QueueBind  *QueueBindConfig
	Consume    *ConsumeConfig
	Publish    *PublishConfig
	Reconnect  *ReconnectConfig
	Ack        *AckConfig
}

type AckConfig struct {
	MaxWorker uint32
	Timeout   time.Duration
}

type ExchangeConfig struct {
	// fanout/direct/topic
	Type       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
}

type QueueConfig struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}

type QueueBindConfig struct {
	RoutingKey string
	NoWait     bool
}

type ConsumeConfig struct {
	AutoAck   bool
	Exclusive bool
	NoWait    bool
}

type PublishConfig struct {
	Mandatory bool
	Immediate bool
	Confirme  bool
	Timeout   time.Duration
}

type ReconnectConfig struct {
	Count     uint
	Delay     time.Duration
	MaxJitter time.Duration
}

// setDefaultConfig
func setDefaultConfig() {

	defaultConfig := map[string]any{
		moduleName + ".vhost":       "/",
		moduleName + ".max_channel": (2 << 15) - 1,

		moduleName + ".exchange.type":       "topic",
		moduleName + ".exchange.durable":    true,
		moduleName + ".exchange.autodelete": false,
		moduleName + ".exchange.internal":   false,
		moduleName + ".exchange.nowait":     false,

		moduleName + ".queue.durable":    true,
		moduleName + ".queue.autodelete": false,
		moduleName + ".queue.exclusive":  false,
		moduleName + ".queue.nowait":     false,

		moduleName + ".bind.routing_key": "",
		moduleName + ".bind.nowait":      false,

		moduleName + ".consume.auto_ack":  false,
		moduleName + ".consume.exclusive": false,
		moduleName + ".consume.nowait":    false,

		moduleName + ".publish.mandatory": false,
		moduleName + ".publish.immediate": false,
		moduleName + ".publish.confirme":  true,
		moduleName + ".publish.timeout":   5 * time.Second,

		moduleName + ".ack.max_worker": 32,
		moduleName + ".ack.timeout":    15 * time.Second,

		moduleName + ".reconnect.count":      10,
		moduleName + ".reconnect.delay":      500 * time.Millisecond,
		moduleName + ".reconnect.max_jitter": 500 * time.Millisecond,
	}

	for k, v := range defaultConfig {
		viper.SetDefault(k, v)
	}
}

func NewRabbitMqConfig() *RabbitMqConfig {
	setDefaultConfig()

	url := viper.GetString(moduleName + ".url")
	vhost := viper.GetString(moduleName + ".vhost")
	channelMax := viper.GetInt(moduleName + ".max_channel")

	exchangeType := viper.GetString(moduleName + ".exchange.type")
	exchangeDurable := viper.GetBool(moduleName + ".exchange.durable")
	exchangeAutoDelete := viper.GetBool(moduleName + ".exchange.autodelete")
	exchangeInternal := viper.GetBool(moduleName + ".exchange.internal")
	exchangeNowait := viper.GetBool(moduleName + ".exchange.nowait")

	queueDurable := viper.GetBool(moduleName + ".queue.durable")
	queueAutoDelete := viper.GetBool(moduleName + ".queue.autodelete")
	queueExclusive := viper.GetBool(moduleName + ".queue.exclusive")
	queueNowait := viper.GetBool(moduleName + ".queue.nowait")

	bindRoutingKey := viper.GetString(moduleName + ".bind.routing_key")
	bindNowait := viper.GetBool(moduleName + ".bind.nowait")

	consumeAutoAck := viper.GetBool(moduleName + ".consume.auto_ack")
	consumeExclusive := viper.GetBool(moduleName + ".consume.exclusive")
	consumeNowait := viper.GetBool(moduleName + ".consume.nowait")

	publishMandatory := viper.GetBool(moduleName + ".publish.mandatory")
	publishImmediate := viper.GetBool(moduleName + ".publish.immediate")
	publishConfirme := viper.GetBool(moduleName + ".publish.confirme")
	publishTimeout := viper.GetDuration(moduleName + ".publish.timeout")

	ackWorker := viper.GetUint32(moduleName + ".ack.max_worker")
	ackTimeout := viper.GetDuration(moduleName + ".ack.timeout")

	reconnectCount := viper.GetUint(moduleName + ".reconnect.count")
	reconnectDelay := viper.GetDuration(moduleName + ".reconnect.delay")
	reconnectMaxJitter := viper.GetDuration(moduleName + ".reconnect.max_jitter")

	return &RabbitMqConfig{
		Url: url,
		AmqpConfig: &amqp.Config{
			Vhost:      vhost,
			ChannelMax: channelMax,
		},
		Exchange: &ExchangeConfig{
			Type:       exchangeType,
			Durable:    exchangeDurable,
			AutoDelete: exchangeAutoDelete,
			Internal:   exchangeInternal,
			NoWait:     exchangeNowait,
		},
		Queue: &QueueConfig{
			Durable:    queueDurable,
			AutoDelete: queueAutoDelete,
			Exclusive:  queueExclusive,
			NoWait:     queueNowait,
		},
		QueueBind: &QueueBindConfig{
			RoutingKey: bindRoutingKey,
			NoWait:     bindNowait,
		},
		Consume: &ConsumeConfig{
			AutoAck:   consumeAutoAck,
			Exclusive: consumeExclusive,
			NoWait:    consumeNowait,
		},
		Publish: &PublishConfig{
			Mandatory: publishMandatory,
			Immediate: publishImmediate,
			Confirme:  publishConfirme,
			Timeout:   publishTimeout,
		},
		Ack: &AckConfig{
			MaxWorker: ackWorker,
			Timeout:   ackTimeout,
		},
		Reconnect: &ReconnectConfig{
			Count:     reconnectCount,
			Delay:     reconnectDelay,
			MaxJitter: reconnectMaxJitter,
		},
	}
}
