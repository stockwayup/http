package controller

import (
	"encoding/json"

	"github.com/stockwayup/http/server/http/response"
	httpSvc "github.com/stockwayup/http/server/http/service"
	"github.com/valyala/fasthttp"
)

// HealthCheck handle actions.
type HealthCheck struct {
	respSvc *httpSvc.Response
}

func NewHealthCheck(respSvc *httpSvc.Response) *HealthCheck {
	return &HealthCheck{respSvc: respSvc}
}

// Handle health check actions.
func (c *HealthCheck) Handle(ctx *fasthttp.RequestCtx) {
	body, err := json.Marshal(response.NewStatus())
	if err != nil {
		c.respSvc.SendInternalError(ctx)

		return
	}

	c.respSvc.SendJSONResponse(ctx, body, fasthttp.StatusOK)
}
