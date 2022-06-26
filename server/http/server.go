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
	cfg         *conf.Config
	http        *controller.HTTP
	healthcheck *controller.HealthCheck
	server      *fasthttp.Server
	logger      *zerolog.Logger
}

func NewServer(
	cfg *conf.Config,
	http *controller.HTTP,
	healthcheck *controller.HealthCheck,
	logger *zerolog.Logger,
) *Server {
	return &Server{
		cfg,
		http,
		healthcheck,
		&fasthttp.Server{},
		logger,
	}
}

// nolint:funlen
func (s *Server) NewRouter() *router.Router {
	r := router.New()
	r.SaveMatchedRoutePath = true

	r.GET("/api/v1/statuses", s.healthcheck.Handle)

	r.POST("/api/v1/users", s.http.Handle)
	r.GET("/api/v1/users/{uid:[0-9]+}/news", s.http.Handle)
	r.GET("/api/v1/users/{uid:[0-9]+}/earnings", s.http.Handle)
	r.GET("/api/v1/users/{uid:[0-9]+}/dividends", s.http.Handle)
	r.GET("/api/v1/users/{uid:[0-9]+}", s.http.Handle)

	r.GET("/api/v1/users/{uid:[0-9]+}/day-prices", s.http.Handle)
	r.GET("/api/v1/users/{uid:[0-9]+}/day-price-periods", s.http.Handle)

	r.GET("/api/v1/users/{uid:[0-9]+}/view-history", s.http.Handle)

	r.POST("/api/v1/refresh-tokens", s.http.Handle)
	r.DELETE("/api/v1/refresh-tokens/{refresh-token}", s.http.Handle)
	r.POST("/api/v1/sessions", s.http.Handle)

	r.GET("/api/v1/confirmation-codes", s.http.Handle)
	r.POST("/api/v1/confirmation-codes/{id}", s.http.Handle)

	r.GET("/api/v1/plans", s.http.Handle)

	r.POST("/api/v1/portfolios", s.http.Handle)
	r.GET("/api/v1/portfolios", s.http.Handle)
	r.GET("/api/v1/portfolios/{pid:[0-9]+}", s.http.Handle)
	r.PATCH("/api/v1/portfolios/{pid:[0-9]+}", s.http.Handle)
	r.DELETE("/api/v1/portfolios/{pid:[0-9]+}", s.http.Handle)

	r.DELETE("/api/v1/portfolios/{pid:[0-9]+}/relationships/securities", s.http.Handle)
	r.POST("/api/v1/portfolios/{pid:[0-9]+}/relationships/securities", s.http.Handle)

	r.POST("/api/v1/portfolios/{pid:[0-9]+}/securities/{sid:[0-9]+}/transactions", s.http.Handle)
	r.GET("/api/v1/portfolios/{pid:[0-9]+}/securities/{sid:[0-9]+}/transactions", s.http.Handle)

	r.PATCH("/api/v1/portfolios/{pid:[0-9]+}/securities/{sid:[0-9]+}/transactions/{tid:[0-9]+}", s.http.Handle)
	r.GET("/api/v1/portfolios/{pid:[0-9]+}/securities/{sid:[0-9]+}/transactions/{tid:[0-9]+}", s.http.Handle)
	r.DELETE("/api/v1/portfolios/{pid:[0-9]+}/securities/{sid:[0-9]+}/transactions/{tid:[0-9]+}", s.http.Handle)
	r.GET("/api/v1/portfolios/{pid:[0-9]+}/securities", s.http.Handle)
	r.GET("/api/v1/portfolios/{pid:[0-9]+}/news", s.http.Handle)
	r.GET("/api/v1/portfolios/{pid:[0-9]+}/earnings", s.http.Handle)
	r.GET("/api/v1/portfolios/{pid:[0-9]+}/dividends", s.http.Handle)
	r.GET("/api/v1/portfolios/{pid:[0-9]+}/day-prices", s.http.Handle)
	r.GET("/api/v1/portfolios/{pid:[0-9]+}/day-price-periods", s.http.Handle)

	r.GET("/api/v1/securities", s.http.Handle)
	r.GET("/api/v1/securities/{sid:[0-9]+}/news", s.http.Handle)
	r.GET("/api/v1/securities/{sid:[0-9]+}/day-prices", s.http.Handle)
	r.GET("/api/v1/securities/{sid:[0-9]+}/day-price-periods", s.http.Handle)

	r.GET("/api/v1/securities/{sid:[0-9]+}/quarterly-balance-sheet", s.http.Handle)
	r.GET("/api/v1/securities/{sid:[0-9]+}/annual-balance-sheet", s.http.Handle)

	r.GET("/api/v1/securities/{sid:[0-9]+}/quarterly-income-statements", s.http.Handle)
	r.GET("/api/v1/securities/{sid:[0-9]+}/annual-income-statements", s.http.Handle)

	r.GET("/api/v1/securities/{sid:[0-9]+}", s.http.Handle)

	r.GET("/api/v1/countries", s.http.Handle)
	r.GET("/api/v1/currencies", s.http.Handle)
	r.GET("/api/v1/sectors", s.http.Handle)
	r.GET("/api/v1/industries", s.http.Handle)
	r.GET("/api/v1/exchanges", s.http.Handle)

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
