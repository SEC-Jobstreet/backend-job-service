package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/SEC-Jobstreet/backend-job-service/utils"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func RunGraphQLServer(ctx context.Context, waitGroup *errgroup.Group, config utils.Config) {
	srv := &http.Server{
		Addr:    config.ListenIP + ":" + config.ListenPort,
		Handler: newRouter(&config),
	}

	var err error
	waitGroup.Go(func() error {
		log.Printf("connect to %s for GraphQL playground", srv.Addr)

		err = srv.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Fatal().Msg("GraphQL API server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown GraphQL API server")

		err = srv.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to shutdown GraphQL API server")
			return err
		}
		log.Info().Msg("GraphQL API server is stopped")
		return nil
	})

	if err != nil {
		log.Fatal().Msg("GraphQL API server failed to serve")
	}
}
