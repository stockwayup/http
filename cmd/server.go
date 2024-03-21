package cmd

import (
	"os"

	"github.com/stockwayup/http/conf"
	"github.com/stockwayup/http/server/http"
	"github.com/stockwayup/http/server/http/controller"
	httpSvc "github.com/stockwayup/http/server/http/service"
	"github.com/stockwayup/http/service"
	"github.com/stockwayup/http/storage/rmq"

	"github.com/rs/zerolog"
	pubsub "github.com/soulgarden/rmq-pubsub"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// nolint: funlen
func NewServerCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Run http server",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			cfg := conf.New()

			logger := zerolog.New(os.Stdout).With().Caller().Logger()

			if cfg.DebugMode {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}

			cmdManager := service.NewManager(&logger)

			ctx, _ := cmdManager.ListenSignal()

			ctx = logger.WithContext(ctx)

			g, ctx := errgroup.WithContext(ctx)

			rmqDialer := rmq.NewDialer(cfg, &logger)
			rmqConn, err := rmqDialer.Dial()
			if err != nil {
				logger.Err(err).Msg("rabbitmq failed to establish connection")
				os.Exit(1)
			}

			defer rmqConn.Close()

			reqBroker := httpSvc.NewRouter()

			pub := pubsub.NewPub(
				rmqConn,
				cfg.RMQ.Queues.Requests,
				pubsub.NewRmq(rmqConn, cfg.RMQ.Queues.Requests, &logger),
				&logger,
			)

			sub := pubsub.NewSub(
				rmqConn,
				httpSvc.NewSubscriber(reqBroker, &logger),
				pubsub.NewRmq(rmqConn, cfg.RMQ.Queues.Responses, &logger),
				cfg.RMQ.Queues.Responses,
				&logger,
			)

			go reqBroker.Start()

			g.Go(func() error {
				return pub.StartPublisher(ctx)
			})

			g.Go(func() error {
				return sub.StartConsumer(ctx)
			})

			respSvc := httpSvc.NewResponse()

			s := http.NewServer(
				cfg,
				controller.NewHTTP(respSvc, pub, reqBroker, &logger),
				controller.NewHealthCheck(respSvc),
				&logger,
			)

			if err != nil {
				logger.Err(err).Msg("router creation failed")
				os.Exit(1)
			}

			g.Go(func() error {
				<-ctx.Done()

				err := s.Shutdown()

				logger.Err(err).Msg("the http server shutdown")

				return err
			})

			g.Go(func() error {
				err := s.Serve(s.NewRouter())

				logger.Err(err).Msg("the http server serve")

				return err
			})

			logger.Err(g.Wait()).Msg("wait goroutines")
		},
	}
}
