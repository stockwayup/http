package controller

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/stockwayup/http/dictionary"

	"github.com/fasthttp/router"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
	pubsub "github.com/soulgarden/rmq-pubsub"
	serverDict "github.com/stockwayup/http/server/http/dictionary"
	httpSvc "github.com/stockwayup/http/server/http/service"
	"github.com/stockwayup/http/storage/rmq/event"
	"github.com/valyala/fasthttp"
)

const (
	authorization = "Authorization"
	bearer        = "Bearer "
	contentType   = "application/octet-stream"
)

type HTTP struct {
	respSvc   *httpSvc.Response
	pub       pubsub.Pub
	reqBroker *httpSvc.Router
	logger    *zerolog.Logger
}

func NewHTTP(
	respSvc *httpSvc.Response,
	pub pubsub.Pub,
	reqBroker *httpSvc.Router,
	logger *zerolog.Logger,
) *HTTP {
	return &HTTP{respSvc: respSvc, pub: pub, reqBroker: reqBroker, logger: logger}
}

func (c *HTTP) Handle(ctx *fasthttp.RequestCtx) {
	msgID := uuid.NewV4().String()
	startedAt := time.Now()

	logCtx := c.logger.With().
		Str("id", msgID).
		Logger().
		WithContext(context.WithValue(ctx, dictionary.ID, msgID))

	defer func() {
		zerolog.Ctx(logCtx).Debug().
			Interface("path", ctx.UserValue(router.MatchedRoutePathParam)).
			Interface("time", time.Since(startedAt).String()).
			Msg("http request processed")
	}()

	req := event.NewHTTPReqByReqCtx(ctx, c.GetAccessTokenFromRequest(ctx))

	bytes, err := req.MarshalMsg(nil)
	if err != nil {
		zerolog.Ctx(logCtx).Err(err).Msg("marshall request")

		c.respSvc.SendInternalError(ctx)
	}

	reqMsg := amqp.Publishing{
		ContentType:  contentType,
		MessageId:    msgID,
		DeliveryMode: amqp.Transient,
		Body:         bytes,
		Timestamp:    time.Now(),
		Expiration:   serverDict.RequestTTL,
	}

	c.pub.Publish(reqMsg)

	ch := c.reqBroker.Subscribe(msgID)

	defer c.reqBroker.Unsubscribe(msgID)

	timeoutCtx, cancel := context.WithTimeout(ctx, serverDict.RequestTimeout)

	defer cancel()

	select {
	case respMsg := <-ch:
		code, err := strconv.Atoi(respMsg.Type)
		if err != nil {
			zerolog.Ctx(logCtx).Err(err).Msg("string to int conversion")

			c.respSvc.SendInternalError(ctx)

			return
		}

		c.respSvc.SendJSONResponse(ctx, respMsg.Body, code)
	case <-timeoutCtx.Done():
		zerolog.Ctx(logCtx).Warn().Msg("timeout")
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
