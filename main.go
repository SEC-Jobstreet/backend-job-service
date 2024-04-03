package main

import (
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/SEC-Jobstreet/backend-job-service/graph"
	"github.com/SEC-Jobstreet/backend-job-service/utils"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	graph.NewDatabaseConnection(&config)

	router := chi.NewRouter()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"ACCEPT", "Authorization", "Content-Type", "X-CSRF-Token", "Origin", "X-Requested-With", "Accept", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		Debug:            true,
		MaxAge:           300,
	}).Handler)
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)
	err = http.ListenAndServe(":4000", router)
	if err != nil {
		panic(err)
	}

	// log.Printf("connect to http://localhost:%s/ for GraphQL playground", config.ListenPort)
	// err = http.ListenAndServe(":"+config.ListenPort, nil)
	// if err != nil {
	// 	log.Fatal().Msg(err.Error())
	// }
}
