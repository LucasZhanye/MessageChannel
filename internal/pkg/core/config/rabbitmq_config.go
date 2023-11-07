package config

import (
	"crypto/tls"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

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
}

// setDefaultConfig
func setDefaultConfig(prefix string) {

	defaultConfig := map[string]any{
		prefix + ".vhost":       "/",
		prefix + ".max_channel": (2 << 15) - 1,

		prefix + ".exchange.type":       "topic",
		prefix + ".exchange.durable":    true,
		prefix + ".exchange.autodelete": false,
		prefix + ".exchange.internal":   false,
		prefix + ".exchange.nowait":     false,

		prefix + ".queue.durable":    true,
		prefix + ".queue.autodelete": false,
		prefix + ".queue.exclusive":  false,
		prefix + ".queue.nowait":     false,

		prefix + ".bind.routing_key": "",
		prefix + ".bind.nowait":      false,

		prefix + ".consume.auto_ack":  false,
		prefix + ".consume.exclusive": false,
		prefix + ".consume.nowait":    false,

		prefix + ".publish.mandatory": false,
		prefix + ".publish.immediate": false,
		prefix + ".publish.confirme":  true,
		prefix + ".publish.timeout":   5 * time.Second,

		prefix + ".ack.max_worker": 32,
		prefix + ".ack.timeout":    15 * time.Second,
	}

	for k, v := range defaultConfig {
		viper.SetDefault(k, v)
	}
}

func NewRabbitMqConfig(module string) *RabbitMqConfig {

	setDefaultConfig(module)

	url := viper.GetString(module + ".url")
	vhost := viper.GetString(module + ".vhost")
	channelMax := viper.GetInt(module + ".max_channel")

	exchangeType := viper.GetString(module + ".exchange.type")
	exchangeDurable := viper.GetBool(module + ".exchange.durable")
	exchangeAutoDelete := viper.GetBool(module + ".exchange.autodelete")
	exchangeInternal := viper.GetBool(module + ".exchange.internal")
	exchangeNowait := viper.GetBool(module + ".exchange.nowait")

	queueDurable := viper.GetBool(module + ".queue.durable")
	queueAutoDelete := viper.GetBool(module + ".queue.autodelete")
	queueExclusive := viper.GetBool(module + ".queue.exclusive")
	queueNowait := viper.GetBool(module + ".queue.nowait")

	bindRoutingKey := viper.GetString(module + ".bind.routing_key")
	bindNowait := viper.GetBool(module + ".bind.nowait")

	consumeAutoAck := viper.GetBool(module + ".consume.auto_ack")
	consumeExclusive := viper.GetBool(module + ".consume.exclusive")
	consumeNowait := viper.GetBool(module + ".consume.nowait")

	publishMandatory := viper.GetBool(module + ".publish.mandatory")
	publishImmediate := viper.GetBool(module + ".publish.immediate")
	publishConfirme := viper.GetBool(module + ".publish.confirme")
	publishTimeout := viper.GetDuration(module + ".publish.timeout")

	ackWorker := viper.GetUint32(module + ".ack.max_worker")
	ackTimeout := viper.GetDuration(module + ".ack.timeout")

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
	}
}
