package controller

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/fasthttp/router"

	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
	pubsub "github.com/soulgarden/rmq-pubsub"
	"github.com/stockwayup/http/server/http/dictionary"
	httpSvc "github.com/stockwayup/http/server/http/service"
	"github.com/stockwayup/http/storage/rmq/event"
	"github.com/streadway/amqp"
	"github.com/valyala/fasthttp"
)

const (
	authorization = "Authorization"
	bearer        = "Bearer "
	contentType   = "application/octet-stream"
)

type HTTP struct {
	respSvc   *httpSvc.Response
	pub       *pubsub.Pub
	reqBroker *httpSvc.Router
	logger    *zerolog.Logger
}

func NewHTTP(
	respSvc *httpSvc.Response,
	pub *pubsub.Pub,
	reqBroker *httpSvc.Router,
	logger *zerolog.Logger,
) *HTTP {
	return &HTTP{respSvc: respSvc, pub: pub, reqBroker: reqBroker, logger: logger}
}

func (c *HTTP) Handle(ctx *fasthttp.RequestCtx) {
	msgID := uuid.NewV4().String()
	startedAt := time.Now()

	logger := c.logger.With().Str("id", msgID).Logger()

	defer func() {
		logger.Debug().
			Str("id", msgID).
			Interface("path", ctx.UserValue(router.MatchedRoutePathParam)).
			Interface("time", time.Since(startedAt).String()).
			Msg("http request processed")
	}()

	req := event.NewHTTPReqByReqCtx(ctx, c.GetAccessTokenFromRequest(ctx))

	bytes, err := req.MarshalMsg(nil)
	if err != nil {
		logger.Err(err).Msg("marshall request")

		c.respSvc.SendInternalError(ctx)
	}

	reqMsg := amqp.Publishing{
		ContentType:  contentType,
		MessageId:    msgID,
		DeliveryMode: amqp.Transient,
		Body:         bytes,
		Timestamp:    time.Now(),
		Expiration:   dictionary.RequestTTL,
	}

	c.pub.Publish(reqMsg)

	ch := c.reqBroker.Subscribe(msgID)

	defer c.reqBroker.Unsubscribe(msgID)

	timeoutCtx, cancel := context.WithTimeout(ctx, dictionary.RequestTimeout)

	defer cancel()

	select {
	case respMsg := <-ch:
		code, err := strconv.Atoi(respMsg.Type)
		if err != nil {
			logger.Err(err).Msg("string to int conversion")

			c.respSvc.SendInternalError(ctx)

			return
		}

		c.respSvc.SendJSONResponse(ctx, respMsg.Body, code)
	case <-timeoutCtx.Done():
		logger.Warn().Msg("timeout")
		c.respSvc.SendTimeoutError(ctx)

		return
	}
}

func (c *HTTP) GetAccessTokenFromRequest(ctx *fasthttp.RequestCtx) string {
	auth := string(ctx.Request.Header.Peek(authorization))
	if auth == "" || !strings.Contains(auth, bearer) {
		return ""
	}

	return strings.Split(auth, bearer)[1]
}
