package http

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type Subscriber struct {
	reqBroker *Router
	logger    *zerolog.Logger
}

func NewSubscriber(reqBroker *Router, logger *zerolog.Logger) *Subscriber {
	return &Subscriber{
		reqBroker: reqBroker,
		logger:    logger,
	}
}

func (s *Subscriber) Process(
	ctx context.Context,
	delivery <-chan amqp.Delivery,
) error {
	for {
		select {
		case msg := <-delivery:
			s.reqBroker.Publish(msg)

			s.logger.Err(msg.Ack(false)).Str("id", msg.MessageId).Msg("ack")
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
