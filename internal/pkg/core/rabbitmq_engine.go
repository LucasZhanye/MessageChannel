package core

import (
	"context"
	"errors"
	"fmt"

	"messagechannel/internal/pkg/core/config"
	"messagechannel/pkg/logger"
	"messagechannel/pkg/pool"
	"messagechannel/pkg/protocol"
	"messagechannel/pkg/safego"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/muesli/cache2go"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Header amqp.Table

type RabbitMqEngine struct {
	config *config.RabbitMqConfig
	log    logger.Log

	conn        *amqp.Connection
	publishPool sync.Pool

	subscriptionManager *SubscriptionManager
	ackCache            *cache2go.CacheTable
	ackChan             chan *protocol.AckRequest
	ackWorkerPool       *pool.WorkerPool

	connected bool
	closeChan chan struct{}
	closed    bool
	lock      sync.Mutex
}

func NewRabbitMqEngine(log logger.Log, subscriptionManager *SubscriptionManager) (*RabbitMqEngine, error) {
	config := config.NewRabbitMqConfig()

	engine := &RabbitMqEngine{
		config:              config,
		log:                 log,
		ackChan:             make(chan *protocol.AckRequest),
		ackWorkerPool:       pool.NewWorkerPool(config.Ack.MaxWorker),
		subscriptionManager: subscriptionManager,
		closeChan:           make(chan struct{}),
	}

	cache := cache2go.Cache("rabbitmq")
	cache.SetOnExpiredCallback(engine.handleAckTimeout)

	engine.ackCache = cache

	return engine, nil
}

func (rabbit *RabbitMqEngine) connect() error {
	rabbit.lock.Lock()
	if rabbit.connected {
		rabbit.lock.Unlock()
		return nil
	}

	defer rabbit.lock.Unlock()

	var conn *amqp.Connection
	var err error

	if rabbit.config.TlsConfig != nil {
		conn, err = amqp.DialTLS(rabbit.config.Url, rabbit.config.TlsConfig)
	} else if rabbit.config.AmqpConfig != nil {
		conn, err = amqp.DialConfig(rabbit.config.Url, *rabbit.config.AmqpConfig)
	} else {
		conn, err = amqp.Dial(rabbit.config.Url)
	}

	if err != nil {
		return fmt.Errorf("cannot connect to rabbitmq: %w", err)
	}

	rabbit.conn = conn
	rabbit.connected = true

	rabbit.log.Info("Connected to rabbitmq engine.")

	return nil
}

func (rabbit *RabbitMqEngine) reconnect() error {
	rabbit.log.Debug("conf = %+v", rabbit.config.Reconnect)
	err := retry.Do(func() error {
		err := rabbit.connect()
		if err == nil {
			return nil
		}

		rabbit.log.Error("Cannot reconnect to rabbitmq,retry...")
		return err
	},
		retry.Attempts(rabbit.config.Reconnect.Count),
		retry.Delay(rabbit.config.Reconnect.Delay),
		retry.MaxJitter(rabbit.config.Reconnect.MaxJitter),
		retry.LastErrorOnly(true))

	return err
}

func (rabbit *RabbitMqEngine) Close() {
	rabbit.lock.Lock()

	if rabbit.closed {
		rabbit.lock.Unlock()
		return
	}

	defer rabbit.lock.Unlock()

	close(rabbit.closeChan)
	rabbit.connected = false
	rabbit.closed = true

	rabbit.log.Info("Closing Rabbitmq Engine...")

	// close rabbitmq connection
	if err := rabbit.conn.Close(); err != nil {
		rabbit.log.Error("Rabbitmq Connection Close error: %v", err)
	}

	rabbit.log.Info("Closed Rabbitmq Engine.")
}

func (rabbit *RabbitMqEngine) Run() error {
	// connect to rabbitmq
	err := rabbit.connect()
	if err != nil {
		return errors.New("Can not connect to rabbitmq engine.")
	}

	err = rabbit.init()
	if err != nil {
		return fmt.Errorf("Can not init rabbitmq engine,err = %w", err)
	}

	// start ack worker
	safego.Execute(rabbit.log, rabbit.ackWorkerPool.Run)

	// handle close
	safego.Execute(rabbit.log, func() {
		for {
			// watch rabbitmq connection not normal close
			notifyClose := rabbit.conn.NotifyClose(make(chan *amqp.Error))

			select {
			case <-rabbit.closeChan:
				return
			case <-notifyClose:
				rabbit.log.Info("Rabbitmq Engine close notify,try to reconnect...")
				rabbit.connected = false
				err := rabbit.reconnect()
				if err != nil {
					rabbit.log.Info("Rabbitmq Engine reconnect %d times fail", rabbit.config.Reconnect.Count)
					// TODO: notify admin
				}
				rabbit.log.Info("Rabbitmq Engine reconnect success")
			}
		}
	})

	return nil
}

func (rabbit *RabbitMqEngine) AcquirePublishChannel() *amqp.Channel {
	channel := rabbit.publishPool.Get()
	if channel == nil {
		c, err := rabbit.conn.Channel()
		if err != nil {
			return nil
		}

		return c
	}

	c := channel.(*amqp.Channel)

	return c
}

func (rabbit *RabbitMqEngine) ReleasePublishChannel(channel *amqp.Channel) {
	rabbit.publishPool.Put(channel)
}

// Publish
func (rabbit *RabbitMqEngine) Publish(publication *protocol.Publication) error {
	if rabbit.closed {
		return errors.New("Rabbitmq Engine is closed")
	}

	if !rabbit.connected {
		return errors.New("Rabbitmq Engine not connected")
	}

	msgId := protocol.GenerateMessageId()

	err := rabbit.publishMessage(publication, msgId, nil)
	if err != nil {
		rabbit.log.Error("client[%s] publish message[id=%s] fail, error = %v", publication.Identifie, msgId, err.Error())
		return err
	}

	rabbit.log.Debug("client[%s] publish message[id=%s] success", publication.Identifie, msgId)
	return nil
}

func (rabbit *RabbitMqEngine) publishMessage(publication *protocol.Publication, msgId string, header Header) error {
	channel := rabbit.AcquirePublishChannel()
	if channel == nil {
		return errors.New("Rabbitmq Engine can not get channel")
	}
	defer rabbit.ReleasePublishChannel(channel)

	var h Header
	if header == nil {
		h = Header{
			"consume_count": 1,
		}
	} else {
		h = header
	}

	var confirmChan chan amqp.Confirmation

	if rabbit.config.Publish.Confirme {
		err := channel.Confirm(false)
		if err != nil {
			rabbit.log.Error("rabbitmq channel can not enter publish confirm mode")
			// cancel confirm mode
			rabbit.config.Publish.Confirme = false
		}

		confirmChan = channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	}

	ctx, cancel := context.WithTimeout(context.Background(), rabbit.config.Publish.Timeout)
	defer cancel()

	routingKey := "#"
	if rk, ok := publication.MetaData[config.ROUTINGKEY]; ok {
		routingKey = rk
	}

	err := channel.PublishWithContext(
		ctx,
		publication.Topic,
		routingKey,
		rabbit.config.Publish.Mandatory,
		rabbit.config.Publish.Immediate,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			MessageId:    msgId,
			AppId:        publication.Identifie,
			Body:         publication.Data,
			Headers:      amqp.Table(h),
		},
	)
	if err != nil {
		return fmt.Errorf("Rabbitmq publish message[id=%s] error: %v", msgId, err)
	}

	if !rabbit.config.Publish.Confirme {
		return nil
	} else {
		// wait until delivery confirmation
		if confimed := <-confirmChan; confimed.Ack {
			return nil
		} else {
			return errors.New("rabbitmq not confirmed for message")
		}
	}
}

// Subscribe
// routingkey equal to topic, Queue name equal to "topicName_groupName", exchange name equal to group name
func (rabbit *RabbitMqEngine) Subscribe(subscription *protocol.Subscription) error {
	if rabbit.closed {
		return errors.New("Rabbitmq Engine is closed")
	}

	if !rabbit.connected {
		return errors.New("Rabbitmq Engine not connected")
	}

	exchangeName := subscription.Topic
	queueName := rabbit.generateQueneName(subscription.Topic, subscription.Group)
	routingkey := "#"
	err := rabbit.prepareSubscribe(exchangeName, queueName, routingkey)
	if err != nil {
		return err
	}

	eChan := make(chan error)

	// start goroutinue to watch message
	safego.Execute(rabbit.log, func() {
		channel, err := rabbit.conn.Channel()
		if err != nil {
			rabbit.log.Error("Rabbitmq subscribe get channel error: %v", err)
			eChan <- err
			return
		}

		defer func() {
			rabbit.log.Info("Exit subscribe goroutinue[topic=%s,group=%s,identifie=%s]", subscription.Topic, subscription.Group, subscription.Identifie)
			if err := channel.Close(); err != nil {
				eChan <- err
			}

			rabbit.subscriptionManager.Remove(subscription.Topic, subscription.Group, subscription.Identifie)
		}()

		notifyClose := channel.NotifyClose(make(chan *amqp.Error))

		// channel.Qos(1, 0, false)

		// consume
		message, err := channel.Consume(
			queueName,
			subscription.Identifie,
			rabbit.config.Consume.AutoAck,
			rabbit.config.Consume.Exclusive,
			false, // The noLocal flag is not supported by RabbitMQ
			rabbit.config.Consume.NoWait,
			nil,
		)

		eChan <- nil
	HandleSubscribe:
		for {
			select {
			case msg := <-message:
				rabbit.handleMessage(msg, subscription)
			case <-notifyClose:
				rabbit.log.Error("Subscribe channel closed informality!")
				break HandleSubscribe

			case <-rabbit.closeChan:
				rabbit.log.Debug("Engine closed so that susbscribe goroutinue exit.")
				break HandleSubscribe

			case <-subscription.ExitChan:
				break HandleSubscribe
			}
		}

	})

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(500*time.Millisecond))
	select {
	case err := <-eChan:
		return err
	case <-ctx.Done():
		return errors.New("Subscribe Timeout")
	}
}

// UnSubscribe
func (rabbit *RabbitMqEngine) UnSubscribe(unsubscription *protocol.Unsubscription) error {
	rabbit.log.Info("UnSubscribe: topic = %s, group = %s, identifie = %s", unsubscription.Topic, unsubscription.Group, unsubscription.Identifie)

	client := rabbit.subscriptionManager.Get(unsubscription.Topic, unsubscription.Group, unsubscription.Identifie)
	if client != nil {
		// client goroutinue exit will remove subscription
		// rabbit.subscriptionManager.Remove(unsubscription.Topic, unsubscription.Group, unsubscription.Identifie)
		client.Info.RemoveSubscription(unsubscription.Topic)
	}
	return nil
}

func (rabbit *RabbitMqEngine) SyncPublish(syncpub *protocol.SyncPublication) (*protocol.Message, error) {
	if rabbit.closed {
		return nil, errors.New("Rabbitmq Engine is closed")
	}

	if !rabbit.connected {
		return nil, errors.New("Rabbitmq Engine not connected")
	}

	// topic name ignore case
	publication := &protocol.Publication{
		Identifie: syncpub.Identifie,
		Topic:     syncpub.Topic,
		Data:      syncpub.Data,
	}
	timeout := syncpub.Timeout
	if timeout <= 0 {

	}
	msgId := protocol.GenerateMessageId()

	channel, err := rabbit.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Rabbitmq connection get channel error:%v", err)
	}

	defer func() {
		// _, err := channel.QueueDelete(msgId, false, false, rabbit.config.Queue.NoWait)
		// if err != nil {
		// 	rabbit.log.Error("Can not delete queue: %s", msgId)
		// }

		if err := channel.Close(); err != nil {
			rabbit.log.Error("rabbitmq channel close error: %v", err)
		}
	}()

	// create temp queue and bind to RESPONSE exchange
	queue, err := channel.QueueDeclare(
		msgId,
		rabbit.config.Queue.Durable,
		rabbit.config.Queue.AutoDelete,
		rabbit.config.Queue.Exclusive,
		rabbit.config.Queue.NoWait,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Rabbitmq declare temp queue error:%v", err)
	}

	err = channel.QueueBind(
		queue.Name,
		msgId,
		config.RESPONSE_EXCHANGE,
		rabbit.config.QueueBind.NoWait,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Rabbitmq bind queue error:%v", err)
	}

	messages, err := channel.Consume(
		msgId,
		msgId,
		rabbit.config.Consume.AutoAck,
		rabbit.config.Consume.Exclusive,
		false, // The noLocal flag is not supported by RabbitMQ
		rabbit.config.Consume.NoWait,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Rabbitmq consume error:%v", err)
	}

	err = rabbit.publishMessage(publication, msgId, nil)
	if err != nil {
		rabbit.log.Error("client[%s] syncpublish message[id=%s] fail, error = %v", publication.Identifie, msgId, err.Error())
		return nil, err
	}

	// start listen
	notifyChannelClose := channel.NotifyClose(make(chan *amqp.Error))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout))
	defer cancel()

	var response *protocol.Message

	rabbit.log.Debug("Waiting response...")
HandleListen:
	select {
	case <-rabbit.closeChan: // connection closed
		err = errors.New("rabbitmq engine connection closed.")
		break HandleListen

	case e := <-notifyChannelClose: // channel closed
		err = fmt.Errorf("rabbitmq engine channel closed, err = %v", e)
		break HandleListen

	case <-ctx.Done(): // timeout
		err = errors.New("sync publish response timeout.")
		break HandleListen

	case msg, ok := <-messages:
		if ok {

			if msg.MessageId != msgId {
				break HandleListen
			}
			response = &protocol.Message{
				Id:        msgId,
				Timestamp: time.Now().Unix(),
				Payload:   msg.Body,
			}
			msg.Ack(false)
		} else {
			msg.Ack(false)
			goto HandleListen
		}
	}

	if response == nil {
		return nil, err
	}

	return response, nil
}

func (rabbit *RabbitMqEngine) SyncPublishReply(syncPubReply *protocol.SyncPublicationReply) error {

	publication := &protocol.Publication{
		Identifie: syncPubReply.Identifie,
		Topic: config.RESPONSE_EXCHANGE,
		Data:  syncPubReply.Data,
		MetaData: map[string]string{
			config.ROUTINGKEY: syncPubReply.MessageId,
		},
	}

	err := rabbit.publishMessage(publication, syncPubReply.MessageId, syncPubReply.Header)
	if err != nil {
		rabbit.log.Error("client[%s] reply message[id=%s] fail, error = %v", publication.Identifie, syncPubReply.MessageId, err.Error())
		return err
	}

	rabbit.log.Debug("client[%s] reply message[id=%s] success", publication.Identifie, syncPubReply.MessageId)
	return nil
}

func (rabbit *RabbitMqEngine) Ack(ackReq *protocol.AckRequest) error {

	task := &pool.Task{
		Fn:    rabbit.handleAck,
		Param: []any{ackReq},
	}

	rabbit.ackWorkerPool.Put(task)

	return nil
}

func (rabbit *RabbitMqEngine) init() error {
	channel, err := rabbit.conn.Channel()
	if err != nil {
		return fmt.Errorf("Rabbitmq connection get channel error:%v", err)
	}

	defer func() {
		if err := channel.Close(); err != nil {
			rabbit.log.Error("rabbitmq channel close error: %v", err)
		}
	}()

	// create dead letter exchange
	err = channel.ExchangeDeclare(
		config.DEAD_LETTER_EXCHANGE,
		rabbit.config.Exchange.Type,
		rabbit.config.Exchange.Durable,
		rabbit.config.Exchange.AutoDelete,
		rabbit.config.Exchange.Internal,
		rabbit.config.Exchange.NoWait,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Rabbitmq declare dead-letter exchange error:%v", err)
	}

	// create dead letter queue
	queue, err := channel.QueueDeclare(
		config.DEAD_LETTER_QUEUE,
		rabbit.config.Queue.Durable,
		rabbit.config.Queue.AutoDelete,
		rabbit.config.Queue.Exclusive,
		rabbit.config.Queue.NoWait,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Rabbitmq declare dead-letter queue error:%v", err)
	}

	err = channel.QueueBind(
		queue.Name,
		config.DEAD_LETTER_KEY,
		config.DEAD_LETTER_EXCHANGE,
		rabbit.config.QueueBind.NoWait,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Rabbitmq bind queue error:%v", err)
	}

	err = channel.ExchangeDeclare(
		config.RESPONSE_EXCHANGE,
		rabbit.config.Exchange.Type,
		rabbit.config.Exchange.Durable,
		rabbit.config.Exchange.AutoDelete,
		rabbit.config.Exchange.Internal,
		rabbit.config.Exchange.NoWait,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Rabbitmq declare response exchange error:%v", err)
	}

	return nil
}

func (rabbit *RabbitMqEngine) generateQueneName(topic, group string) string {
	return topic + "_" + group
}

func (rabbit *RabbitMqEngine) prepareSubscribe(exchangeName, queueName, routingkey string) error {
	channel, err := rabbit.conn.Channel()
	if err != nil {
		return fmt.Errorf("Rabbitmq connection get channel error:%v", err)
	}

	defer func() {
		if err := channel.Close(); err != nil {
			rabbit.log.Error("rabbitmq channel close error: %v", err)
		}
	}()

	// declare exchange
	err = channel.ExchangeDeclare(
		exchangeName,
		rabbit.config.Exchange.Type,
		rabbit.config.Exchange.Durable,
		rabbit.config.Exchange.AutoDelete,
		rabbit.config.Exchange.Internal,
		rabbit.config.Exchange.NoWait,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Rabbitmq declare exchange error:%v", err)
	}

	// create queue
	queue, err := channel.QueueDeclare(
		queueName,
		rabbit.config.Queue.Durable,
		rabbit.config.Queue.AutoDelete,
		rabbit.config.Queue.Exclusive,
		rabbit.config.Queue.NoWait,
		amqp.Table{
			"x-dead-letter-exchange":    config.DEAD_LETTER_EXCHANGE,
			"x-dead-letter-routing-key": config.DEAD_LETTER_KEY,
		},
	)
	if err != nil {
		return fmt.Errorf("Rabbitmq declare queue error:%v", err)
	}

	err = channel.QueueBind(
		queue.Name,
		routingkey,
		exchangeName,
		rabbit.config.QueueBind.NoWait,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Rabbitmq bind queue error:%v", err)
	}

	return nil
}

func (rabbit *RabbitMqEngine) handleMessage(message amqp.Delivery, subscription *protocol.Subscription) {
	rabbit.log.Info("recevie msg = %v, header = %+v", string(message.Body), message.Headers)
	// rabbit.log.Debug("topic = %s, group = %s, id = %s", subscription.Topic, subscription.Group, subscription.Identifie)

	client := rabbit.subscriptionManager.Get(subscription.Topic, subscription.Group, subscription.Identifie)
	if client != nil {
		response := protocol.AcquireResponse(protocol.MESSAGE_RESPONSE_CMDID)
		defer protocol.ReleaseResponse(response)

		sendMsg := protocol.NewMessage(message.MessageId, message.Timestamp.Unix(), message.Body)
		response = response.SetMessageResponse(sendMsg)

		err := client.Send(response.ToJsonBytes())
		if err != nil {
			// TODO: handle error, maybe retry
			rabbit.log.Error("send data to client[%s] fail", subscription.Identifie)
			return
		}

		rabbit.log.Debug("set to cache: %s, %s, %s", message.MessageId, message.RoutingKey, message.Exchange)
		// set group name to message header
		message.Headers["group"] = subscription.Group
		rabbit.ackCache.Add(message.MessageId, rabbit.config.Ack.Timeout, message)
	} else {
		// TODO: Maybe should send to error queue
		message.Ack(false)
	}
}

func (rabbit *RabbitMqEngine) handleAck(params ...any) {

	ackReq := params[0].(*protocol.AckRequest)

	cacheItem, err := rabbit.ackCache.Value(ackReq.MessageId)
	if err != nil {
		// note: means ack timeout, ignore it
		return
	}

	message := cacheItem.Data().(amqp.Delivery)

	result := ackReq.Result

	if result == protocol.ACKRESULT_ACK {
		message.Ack(false)
	} else {
		message.Nack(false, false) // send to dead-letter exchange
	}

	// remove from cache
	rabbit.ackCache.Delete(ackReq.MessageId)
}

func (rabbit *RabbitMqEngine) handleAckTimeout(item *cache2go.CacheItem) {
	rabbit.log.Debug("ack timeout, ack message id = %v", item.Key())

	message := item.Data().(amqp.Delivery)

	groupName, ok := message.Headers["group"].(string)
	if !ok {
		message.Nack(false, false) // send to dead-letter exchange
		return
	}
	consumerCount := rabbit.subscriptionManager.GetSubscriptionCount(message.Exchange, groupName)
	if consumerCount > 1 {
		consumeTimes := int(message.Headers["consume_count"].(int32))
		// find next consumer
		if consumeTimes < consumerCount {
			message.Headers["consume_count"] = consumeTimes + 1
			// republish
			publication := &protocol.Publication{
				Identifie: message.AppId,
				Topic:     message.Exchange,
				Data:      message.Body,
			}
			err := rabbit.publishMessage(publication, message.MessageId, Header(message.Headers))
			if err != nil {
				rabbit.log.Error("client[%s] transfer message[id=%s] to next client fail, error = %v", publication.Identifie, message.MessageId, err.Error())
				message.Nack(false, false)
			} else {
				message.Ack(false)
			}
		} else {
			message.Nack(false, false)
		}
	} else {
		message.Nack(false, false) // send to dead-letter exchange
	}
}
