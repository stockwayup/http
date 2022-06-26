package http

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
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
) {
	for {
		select {
		case msg := <-delivery:
			s.reqBroker.Publish(msg)

			s.logger.Err(msg.Ack(false)).Str("id", msg.MessageId).Msg("ack")
		case <-ctx.Done():
			return
		}
	}
}
