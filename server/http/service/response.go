package http

import (
	"encoding/json"

	"github.com/stockwayup/http/server/http/response"

	"github.com/valyala/fasthttp"
)

const contentType = "application/vnd.api+json; charset=utf-8"

type Response struct {
}

func NewResponse() *Response {
	return &Response{}
}

func (s *Response) SendJSONResponse(ctx *fasthttp.RequestCtx, body []byte, code int) {
	s.sendJSONResponseWithBody(ctx, body, code)
}

func (s *Response) SendInternalError(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
}

func (s *Response) SendNoContent(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNoContent)
}

func (s *Response) SendBadRequest(ctx *fasthttp.RequestCtx, body interface{}) {
	s.marshallAndSend(ctx, body, fasthttp.StatusBadRequest)
}

func (s *Response) SendUnauthorized(ctx *fasthttp.RequestCtx) {
	s.marshallAndSend(ctx, response.Errors{
		Errors: []response.Error{response.UnauthorizedError},
	}, fasthttp.StatusUnauthorized)
}

func (s *Response) SendTimeoutError(ctx *fasthttp.RequestCtx) {
	s.marshallAndSend(ctx, response.Errors{
		Errors: []response.Error{response.TimeoutError},
	}, fasthttp.StatusRequestTimeout)
}

func (s *Response) SendNotFound(ctx *fasthttp.RequestCtx) {
	s.marshallAndSend(ctx, response.Errors{
		Errors: []response.Error{response.NotFoundError},
	}, fasthttp.StatusNotFound)
}

func (s *Response) SendAccessDenied(ctx *fasthttp.RequestCtx) {
	s.marshallAndSend(ctx, response.Errors{
		Errors: []response.Error{response.ForbiddenError},
	}, fasthttp.StatusForbidden)
}

func (s *Response) sendJSONResponseWithBody(ctx *fasthttp.RequestCtx, body []byte, code int) {
	ctx.SetContentType(contentType)
	ctx.SetStatusCode(code)
	ctx.Response.SetBody(body)
}

func (s *Response) marshallAndSend(ctx *fasthttp.RequestCtx, body interface{}, code int) {
	marshaled, err := json.Marshal(body)
	if err != nil {
		s.SendInternalError(ctx)

		return
	}

	s.sendJSONResponseWithBody(ctx, marshaled, code)
}
