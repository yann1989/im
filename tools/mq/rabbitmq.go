// Author: yann
// Date: 2020/5/26 9:53 上午
// Desc:

package mq

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sync"
	chat "yann-chat"
)

type Config struct {
	AmqpAddr     string
	ExchangeName string
}

const (
	EXCHANGE_KIND_FANOUT = "fanout"
	CONTENT_TYPE         = "text/plain"
)

var (
	once         sync.Once
	amqpProducer *amqp.Connection
	amqpConsumer *amqp.Connection
	chProducer   *amqp.Channel
	chConsumer   *amqp.Channel
	exchangeName string
)

//初始化参数
func (c *Config) initParam() {
	if c.ExchangeName == "" {
		logrus.Errorf("启动失败,请初始化mq交换机名称")
		panic("请初始化mq交换机名称")
	}
	if c.AmqpAddr == "" {
		c.AmqpAddr = "amqp://guest:guest@127.0.0.1:5672/"
	}
}

func (c *Config) Start(ctx context.Context, yannChat *chat.YannChat) error {
	once.Do(func() {
		c.initParam()
		var err error
		amqpProducer, err = amqp.Dial(c.AmqpAddr)
		if err != nil {
			logrus.Errorf("amqp 初始化失败, 失败原因:%s", err.Error())
			panic("amqp 初始化失败, 失败原因:%s" + err.Error())
		}
		amqpConsumer, err = amqp.Dial(c.AmqpAddr)
		if err != nil {
			logrus.Errorf("amqp 初始化失败, 失败原因:%s", err.Error())
			panic("amqp 初始化失败, 失败原因:%s" + err.Error())
		}
		chProducer, err = amqpProducer.Channel()
		if err != nil {
			logrus.Errorf("amqp 初始化channel失败, 失败原因:%s", err.Error())
			panic("amqp 初始化channel失败, 失败原因:%s" + err.Error())
		}
		chConsumer, err = amqpProducer.Channel()
		if err != nil {
			logrus.Errorf("amqp 初始化channel失败, 失败原因:%s", err.Error())
			panic("amqp 初始化channel失败, 失败原因:%s" + err.Error())
		}
		//申明生产者交换机
		err = chProducer.ExchangeDeclare(
			c.ExchangeName,       //交换机名称
			EXCHANGE_KIND_FANOUT, //交换机类型广播
			true,                 //durable  持久化
			false,                //autodelete
			false,
			false,
			nil,
		)
		if err != nil {
			logrus.Errorf("amqp 初始化chProducer Exchange失败, 失败原因:%s", err.Error())
			panic("amqp 初始化Exchange失败, 失败原因:%s" + err.Error())
		}
		//申明消费者交换机
		err = chConsumer.ExchangeDeclare(
			c.ExchangeName,       //交换机名称
			EXCHANGE_KIND_FANOUT, //交换机类型广播
			true,                 //durable
			false,                //autodelete
			false,
			false,
			nil,
		)
		if err != nil {
			logrus.Errorf("amqp 初始化chConsumer Exchange失败, 失败原因:%s", err.Error())
			panic("amqp 初始化Exchange失败, 失败原因:%s" + err.Error())
		}
		exchangeName = c.ExchangeName
	})
	return nil
}

func Broadcast(msg []byte) error {
	return chProducer.Publish(
		exchangeName, //exchange
		"",           //routing key
		false,
		false,
		amqp.Publishing{
			ContentType: CONTENT_TYPE,
			Body:        msg,
		})
}

func StartConsume() (<-chan amqp.Delivery, error) {
	queue, err := chConsumer.QueueDeclare(
		"",
		true, //durable
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		logrus.Errorf("amqp chConsumer 声明Queue失败, 失败原因:%s", err.Error())
		panic("amqp chConsumer 声明Queue失败, 失败原因:%s" + err.Error())
	}
	err = chConsumer.QueueBind(
		queue.Name,
		"",
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		logrus.Errorf("amqp chConsumer 绑定Queue失败, 失败原因:%s", err.Error())
		panic("amqp chConsumer 绑定Queue失败Queue失败, 失败原因:%s" + err.Error())
	}

	return chConsumer.Consume(
		queue.Name,
		"",
		true, //Auto Ack
		false,
		false,
		false,
		nil,
	)
}

func (c *Config) Stop(ctx context.Context) (err error) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("mq 关闭异常")
		}
	}()
	amqpProducer.Close()
	amqpConsumer.Close()
	chProducer.Close()
	chConsumer.Close()
	return
}
