package entity

import amqp "github.com/rabbitmq/amqp091-go"

type Sub struct {
	ID string
	Ch chan amqp.Delivery
}
