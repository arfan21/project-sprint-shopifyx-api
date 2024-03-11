package api

import (
	"context"

	"github.com/arfan21/shopifyx-api/config"
	"github.com/arfan21/shopifyx-api/internal/server"
	dbpostgres "github.com/arfan21/shopifyx-api/pkg/db/postgres"
	"github.com/arfan21/shopifyx-api/pkg/logger"
	"github.com/arfan21/shopifyx-api/pkg/telemetry"
	"github.com/urfave/cli/v2"
)

func Serve() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Run the API server",
		Action: func(c *cli.Context) error {
			_, err := config.LoadConfig()
			if err != nil {
				return err
			}

			_, err = config.ParseConfig(config.GetViper())
			if err != nil {
				return err
			}

			if config.Get().Otel.Enabled {
				tracerShutdown, err := telemetry.InitTracer()
				if err != nil {
					return err
				}
				defer tracerShutdown(context.Background())

				logShutdown, err := telemetry.InitLogs()
				if err != nil {
					return err
				}

				defer logShutdown(context.Background())

				metricShutdown, err := telemetry.InitMetric()
				if err != nil {
					return err
				}

				defer metricShutdown(context.Background())

				// this log called for initialize hook
				logger.Log(context.Background()).Info().Msgf("tracing enabled with service name: %s", config.Get().Service.Name)
			} else {
				logger.Log(context.Background()).Warn().Msg("tracing disabled")
			}

			db, err := dbpostgres.NewPgx()
			if err != nil {
				return err
			}

			server := server.New(
				db,
			)
			return server.Run()
		},
	}

}
