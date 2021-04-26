package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"sync"
	"time"
)

type RMQClient struct {
	publishCh         *amqp.Channel
	consumeCh         *amqp.Channel
	conn              *amqp.Connection
	uri               *amqp.URI
	mutex             *sync.Mutex
	connCloseError    chan *amqp.Error
	pubishChCloseErr  chan *amqp.Error
	consumeChCloseErr chan *amqp.Error
}

var isInit = false

func New(uri *amqp.URI) *RMQClient {
	if !isInit {
		isInit = true
	}

	rmq := &RMQClient{
		uri:               uri,
		mutex:             new(sync.Mutex),
		connCloseError:    make(chan *amqp.Error),
		pubishChCloseErr:  make(chan *amqp.Error),
		consumeChCloseErr: make(chan *amqp.Error),
	}

	rmq.connectToRabbit()
	rmq.conn.NotifyClose(rmq.connCloseError)
	rmq.publishCh.NotifyClose(rmq.pubishChCloseErr)
	rmq.consumeCh.NotifyClose(rmq.consumeChCloseErr)

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
	log.Println("rmq: Connect success")
}

func (rmq *RMQClient) reConnector() {
	for {
		select {
		case connError := <-rmq.connCloseError:
			log.Println("rmq: lost connection " + connError.Error())
			log.Println("rmq: begin to reConnect")
			rmq.connectToRabbit()
			rmq.connCloseError = make(chan *amqp.Error)
			rmq.conn.NotifyClose(rmq.connCloseError)
		case publishChErr := <-rmq.pubishChCloseErr:
			log.Println("rmq: publishChErr", publishChErr.Error())
			if rmq.conn != nil {
				rmq.publishCh, _ = rmq.conn.Channel()
			}
			rmq.pubishChCloseErr = make(chan *amqp.Error)
			rmq.publishCh.NotifyClose(rmq.pubishChCloseErr)
		case consumeChErr := <-rmq.consumeChCloseErr:
			log.Println("rmq: consumeChErr", consumeChErr.Error())
			if rmq.conn != nil {
				rmq.consumeCh, _ = rmq.conn.Channel()
			}
			rmq.consumeChCloseErr = make(chan *amqp.Error)
			rmq.publishCh.NotifyClose(rmq.consumeChCloseErr)
		}
	}
}
