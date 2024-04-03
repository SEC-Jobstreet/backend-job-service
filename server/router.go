package server

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/SEC-Jobstreet/backend-job-service/graph"
	"github.com/SEC-Jobstreet/backend-job-service/utils"
	"github.com/gin-gonic/gin"
)

func newRouter(config *utils.Config) *gin.Engine {
	router := gin.Default()

	authRoutes := router.Group("/api/v1").Use(IsAuthorizedJWT(config))

	authRoutes.POST("/", graphqlHandler())
	authRoutes.GET("/query", playgroundHandler())

	return router
}

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
