package entity

import "github.com/streadway/amqp"

type Sub struct {
	ID string
	Ch chan amqp.Delivery
}
