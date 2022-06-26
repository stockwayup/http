package rmq

import (
	"time"

	"github.com/stockwayup/http/conf"

	"github.com/isayme/go-amqp-reconnect/rabbitmq"
	"github.com/rs/zerolog"
)

type Dialer struct {
	cfg    *conf.Config
	logger *zerolog.Logger
}

func NewDialer(cfg *conf.Config, logger *zerolog.Logger) *Dialer {
	return &Dialer{cfg: cfg, logger: logger}
}

func (rmq *Dialer) Dial() (*rabbitmq.Connection, error) {
	rabbitmq.Debug = rmq.cfg.DebugMode

	var conn *rabbitmq.Connection

	var err error

	attempts := 30

	for attempts > 0 {
		conn, err = rabbitmq.Dial(
			"amqp://" + rmq.cfg.RMQ.User + ":" + rmq.cfg.RMQ.Password + "@" + rmq.cfg.RMQ.Host + ":" + rmq.cfg.RMQ.Port,
		)
		rmq.logger.Err(err).Int("attempt", attempts).Msg("dial rabbitmq")

		if err != nil {
			time.Sleep(time.Second)
			attempts--

			continue
		}

		break
	}

	return conn, err
}
