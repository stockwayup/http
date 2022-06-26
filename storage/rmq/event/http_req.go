package event

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

//go:generate msgp

type HTTPReq struct {
	UUID        string                 `msgp:"uuid"`
	Type        string                 `msgp:"type"`
	Body        []byte                 `msgp:"body"`
	AccessToken string                 `msgp:"access_token"`
	Method      string                 `msgp:"method"`
	UserValues  map[string]interface{} `msgp:"user_values"`
	Uri         URI                    `msgp:"uri"` // nolint: golint,revive,stylecheck
}

func (r HTTPReq) URI() URI {
	return r.Uri
}

type URI struct {
	PathOriginal []byte `msgp:"path_original"`
	Scheme       []byte `msgp:"scheme"`
	Path         []byte `msgp:"path"`
	QueryString  []byte `msgp:"query_string"`
	Hash         []byte `msgp:"hash"`
	Host         []byte `msgp:"host"`
	Args         Args   `msgp:"query_args"`
}

func (u URI) QueryArgs() Args {
	return u.Args
}

type Args struct {
	Val map[string][]byte
}

func (a Args) Peek(key string) []byte {
	val, ok := a.Val[key]
	if !ok {
		return nil
	}

	return val
}

func NewHTTPReqByReqCtx(ctx *fasthttp.RequestCtx, accessToken string) HTTPReq {
	// nolint: forcetypeassert
	req := HTTPReq{
		Type:        ctx.UserValue(router.MatchedRoutePathParam).(string),
		Body:        ctx.Request.Body(),
		AccessToken: accessToken,
		Method:      string(ctx.Request.Header.Method()),
		UserValues:  map[string]interface{}{},
		Uri: URI{
			PathOriginal: ctx.Request.URI().PathOriginal(),
			Scheme:       ctx.Request.URI().Scheme(),
			Path:         ctx.Request.URI().Path(),
			QueryString:  ctx.Request.URI().QueryString(),
			Hash:         ctx.Request.URI().Hash(),
			Host:         ctx.Request.URI().Host(),
			Args:         Args{Val: map[string][]byte{}},
		},
	}

	ctx.VisitUserValues(func(key []byte, val interface{}) {
		req.UserValues[string(key)] = val
	})

	ctx.Request.URI().QueryArgs().VisitAll(func(key, value []byte) {
		req.Uri.Args.Val[string(key)] = value
	})

	return req
}
