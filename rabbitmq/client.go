package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"sync"
	"time"
)

type RMQClient struct {
	publishCh  *amqp.Channel
	consumeCh  *amqp.Channel
	conn       *amqp.Connection
	uri        *amqp.URI
	mutex      *sync.Mutex
	closeError chan *amqp.Error
}

var isInit = false

func New(uri *amqp.URI) *RMQClient {
	if !isInit {
		isInit = true
	}

	rmq := &RMQClient{
		uri:        uri,
		mutex:      new(sync.Mutex),
		closeError: make(chan *amqp.Error),
	}

	rmq.connectToRabbit()
	rmq.conn.NotifyClose(rmq.closeError)
	go rmq.reConnector()
	return rmq
}

func (rmq *RMQClient) connectToRabbit() {
	for {
		conn, err := amqp.Dial(rmq.uri.String())
		if err != nil {
			log.Println("rmq: failed to dial RabbitMQ:" + err.Error())
			time.Sleep(time.Second)
			continue
		}

		rmq.conn = conn

		rmq.publishCh, err = conn.Channel()
		if err != nil {
			log.Println("rmq: failed to open channel:" + err.Error())
			time.Sleep(time.Second)
			continue
		}

		rmq.consumeCh, err = conn.Channel()
		if err != nil {
			log.Println("rmq: failed to open channel:" + err.Error())
			time.Sleep(time.Second)
			continue
		}

		break
	}
}

func (rmq *RMQClient) reConnector() {
	for {
		rabbitErr := <-rmq.closeError
		log.Println("rmq: lost connection " + rabbitErr.Error())
		log.Println("rmq: begin to reConnect")
		rmq.publishCh.Close()
		rmq.consumeCh.Close()
		rmq.conn.Close()
		rmq.connectToRabbit()
		rmq.closeError = make(chan *amqp.Error)
		rmq.conn.NotifyClose(rmq.closeError)
		log.Println("rmq: reConnect success")
	}
}
