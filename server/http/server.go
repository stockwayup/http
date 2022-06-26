package http

import (
	cors "github.com/AdhityaRamadhanus/fasthttpcors"
	"github.com/fasthttp/router"
	"github.com/rs/zerolog"
	"github.com/stockwayup/http/conf"
	"github.com/stockwayup/http/server/http/controller"
	"github.com/valyala/fasthttp"
)

const maxAge = 60 * 60

type Server struct {
	cfg        *conf.Config
	controller *controller.HTTP
	server     *fasthttp.Server
	logger     *zerolog.Logger
}

func NewServer(
	cfg *conf.Config,
	controller *controller.HTTP,
	logger *zerolog.Logger,
) *Server {
	return &Server{
		cfg,
		controller,
		&fasthttp.Server{},
		logger,
	}
}

// nolint:funlen
func (s *Server) NewRouter() *router.Router {
	r := router.New()
	r.SaveMatchedRoutePath = true

	r.GET("/api/v1/statuses", s.controller.Handle)

	r.POST("/api/v1/users", s.controller.Handle)
	r.GET("/api/v1/users/{uid:[0-9]+}/news", s.controller.Handle)
	r.GET("/api/v1/users/{uid:[0-9]+}/earnings", s.controller.Handle)
	r.GET("/api/v1/users/{uid:[0-9]+}/dividends", s.controller.Handle)
	r.GET("/api/v1/users/{uid:[0-9]+}", s.controller.Handle)

	r.GET("/api/v1/users/{uid:[0-9]+}/day-prices", s.controller.Handle)
	r.GET("/api/v1/users/{uid:[0-9]+}/day-price-periods", s.controller.Handle)

	r.GET("/api/v1/users/{uid:[0-9]+}/view-history", s.controller.Handle)

	r.POST("/api/v1/refresh-tokens", s.controller.Handle)
	r.DELETE("/api/v1/refresh-tokens/{refresh-token}", s.controller.Handle)
	r.POST("/api/v1/sessions", s.controller.Handle)

	r.GET("/api/v1/confirmation-codes", s.controller.Handle)
	r.POST("/api/v1/confirmation-codes/{id}", s.controller.Handle)

	r.GET("/api/v1/plans", s.controller.Handle)

	r.POST("/api/v1/portfolios", s.controller.Handle)
	r.GET("/api/v1/portfolios", s.controller.Handle)
	r.GET("/api/v1/portfolios/{pid:[0-9]+}", s.controller.Handle)
	r.PATCH("/api/v1/portfolios/{pid:[0-9]+}", s.controller.Handle)
	r.DELETE("/api/v1/portfolios/{pid:[0-9]+}", s.controller.Handle)

	r.DELETE("/api/v1/portfolios/{pid:[0-9]+}/relationships/securities", s.controller.Handle)
	r.POST("/api/v1/portfolios/{pid:[0-9]+}/relationships/securities", s.controller.Handle)

	r.POST("/api/v1/portfolios/{pid:[0-9]+}/securities/{sid:[0-9]+}/transactions", s.controller.Handle)
	r.GET("/api/v1/portfolios/{pid:[0-9]+}/securities/{sid:[0-9]+}/transactions", s.controller.Handle)

	r.PATCH("/api/v1/portfolios/{pid:[0-9]+}/securities/{sid:[0-9]+}/transactions/{tid:[0-9]+}", s.controller.Handle)
	r.GET("/api/v1/portfolios/{pid:[0-9]+}/securities/{sid:[0-9]+}/transactions/{tid:[0-9]+}", s.controller.Handle)
	r.DELETE("/api/v1/portfolios/{pid:[0-9]+}/securities/{sid:[0-9]+}/transactions/{tid:[0-9]+}", s.controller.Handle)
	r.GET("/api/v1/portfolios/{pid:[0-9]+}/securities", s.controller.Handle)
	r.GET("/api/v1/portfolios/{pid:[0-9]+}/news", s.controller.Handle)
	r.GET("/api/v1/portfolios/{pid:[0-9]+}/earnings", s.controller.Handle)
	r.GET("/api/v1/portfolios/{pid:[0-9]+}/dividends", s.controller.Handle)
	r.GET("/api/v1/portfolios/{pid:[0-9]+}/day-prices", s.controller.Handle)
	r.GET("/api/v1/portfolios/{pid:[0-9]+}/day-price-periods", s.controller.Handle)

	r.GET("/api/v1/securities", s.controller.Handle)
	r.GET("/api/v1/securities/{sid:[0-9]+}/news", s.controller.Handle)
	r.GET("/api/v1/securities/{sid:[0-9]+}/day-prices", s.controller.Handle)
	r.GET("/api/v1/securities/{sid:[0-9]+}/day-price-periods", s.controller.Handle)

	r.GET("/api/v1/securities/{sid:[0-9]+}/quarterly-balance-sheet", s.controller.Handle)
	r.GET("/api/v1/securities/{sid:[0-9]+}/annual-balance-sheet", s.controller.Handle)

	r.GET("/api/v1/securities/{sid:[0-9]+}/quarterly-income-statements", s.controller.Handle)
	r.GET("/api/v1/securities/{sid:[0-9]+}/annual-income-statements", s.controller.Handle)

	r.GET("/api/v1/securities/{sid:[0-9]+}", s.controller.Handle)

	r.GET("/api/v1/countries", s.controller.Handle)
	r.GET("/api/v1/currencies", s.controller.Handle)
	r.GET("/api/v1/sectors", s.controller.Handle)
	r.GET("/api/v1/industries", s.controller.Handle)
	r.GET("/api/v1/exchanges", s.controller.Handle)

	return r
}

func (s *Server) Serve(r *router.Router) error {
	s.server.Handler = r.Handler

	if s.cfg.EnableCors {
		withCors := cors.NewCorsHandler(cors.Options{
			AllowedOrigins: []string{
				"http://127.0.0.1",
				"http://localhost",
				"http://127.0.0.1:8080",
				"http://localhost:8080",
				"http://127.0.0.1:8081",
				"http://localhost:8081",
				"https://dev.stockwayup.com",
				"https://stockwayup.com",
			},
			AllowedHeaders:   []string{"origin", "accept", "content-type", "authorization"},
			AllowedMethods:   []string{"GET", "POST", "DELETE", "PATCH"},
			AllowCredentials: false,
			AllowMaxAge:      maxAge,
			Debug:            false,
		})

		s.server.Handler = withCors.CorsMiddleware(r.Handler)
	}

	s.logger.Info().Str("port", s.cfg.ListenPort).Msg("start listening http connections")

	err := s.server.ListenAndServe(":" + s.cfg.ListenPort)

	s.logger.Err(err).Msg("server down")

	return err
}

func (s *Server) Shutdown() error {
	err := s.server.Shutdown()

	s.logger.Err(err).Msg("the http server shutdown")

	return err
}
