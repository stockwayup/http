package http

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stockwayup/http/server/http/entity"
)

const eventChSize = 1

type Router struct {
	subscribers map[string]chan amqp.Delivery
	subCh       chan entity.Sub
	unsubCh     chan string
	publishCh   chan amqp.Delivery
}

func NewRouter() *Router {
	return &Router{
		subscribers: make(map[string]chan amqp.Delivery),
		subCh:       make(chan entity.Sub),
		unsubCh:     make(chan string),
		publishCh:   make(chan amqp.Delivery),
	}
}

func (s *Router) Start() {
	for {
		select {
		case ch := <-s.subCh:
			s.subscribers[ch.ID] = ch.Ch
		case id := <-s.unsubCh:
			if ch, ok := s.subscribers[id]; !ok {
				delete(s.subscribers, id)
				close(ch)
			}

		case msg := <-s.publishCh:
			if subscriber, ok := s.subscribers[msg.MessageId]; ok {
				subscriber <- msg
			}
		}
	}
}

func (s *Router) Subscribe(id string) chan amqp.Delivery {
	msgCh := make(chan amqp.Delivery, eventChSize)

	s.subCh <- entity.Sub{
		ID: id,
		Ch: msgCh,
	}

	return msgCh
}

func (s *Router) Unsubscribe(id string) {
	s.unsubCh <- id
}

func (s *Router) Publish(msg amqp.Delivery) {
	s.publishCh <- msg
}
