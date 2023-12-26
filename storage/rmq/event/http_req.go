package event

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

//go:generate msgp

type HTTPReq struct {
	Type        string                 `msg:"type"`
	AccessToken string                 `msg:"access_token"`
	Method      string                 `msg:"method"`
	UserValues  map[string]interface{} `msg:"user_values"`
	Uri         URI                    `msg:"uri"` //nolint: golint,revive,stylecheck
	Body        []byte                 `msg:"body"`
}

func (r HTTPReq) URI() URI {
	return r.Uri
}

type URI struct {
	PathOriginal []byte `msg:"path_original"`
	Scheme       []byte `msg:"scheme"`
	Path         []byte `msg:"path"`
	QueryString  []byte `msg:"query_string"`
	Hash         []byte `msg:"hash"`
	Host         []byte `msg:"host"`
	Args         Args   `msg:"args"`
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
	//nolint: forcetypeassert
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
