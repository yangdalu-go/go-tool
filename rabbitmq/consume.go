package rabbitmq

import "github.com/streadway/amqp"

func (rmq *RMQClient) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return rmq.consumeCh.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

func (rmq *RMQClient) QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (int, error) {
	return rmq.consumeCh.QueueDelete(name, ifUnused, ifEmpty, noWait)
}

func (rmq *RMQClient) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	return rmq.consumeCh.QueueBind(name, key, exchange, noWait, args)
}

func (rmq *RMQClient) QueueUnbind(name, key, exchange string, args amqp.Table) error {
	return rmq.consumeCh.QueueUnbind(name, key, exchange, args)
}

func (rmq *RMQClient) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return rmq.consumeCh.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
}

func (rmq *RMQClient) ExchangeBind(destination, key, source string, noWait bool, args amqp.Table) error {
	return rmq.consumeCh.ExchangeBind(destination, key, source, noWait, args)
}

func (rmq *RMQClient) ExchangeUnbind(destination, key, source string, noWait bool, args amqp.Table) error {
	return rmq.consumeCh.ExchangeUnbind(destination, key, source, noWait, args)
}

func (rmq *RMQClient) ExchangeDelete(name string, ifUnused, noWait bool) error {
	return rmq.consumeCh.ExchangeDelete(name, ifUnused, noWait)
}

func (rmq *RMQClient) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return rmq.consumeCh.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}
