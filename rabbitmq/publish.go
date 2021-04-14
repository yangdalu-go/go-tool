package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

func (rmq *RMQClient) PublishToExchange(exchange, key string, headers amqp.Table, body interface{}) error {
	jsonStr, err := json.Marshal(body)
	if err != nil {
		return err
	}

	pub := amqp.Publishing{
		Headers:     headers,
		ContentType: "text/json",
		Body:        jsonStr,
	}

	return rmq.publishCh.Publish(exchange, key, false, false, pub)
}
