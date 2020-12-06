package spring

import (
	"fmt"
	"os"

	"github.com/99designs/gqlgen/graphql"
)

func (s *Spring) RunServer(defaultPort uint, server graphql.ExecutableSchema) {
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", defaultPort)
	}

	s.initializeIoCHandlers()
	s.initializeLog()
	ginEngine := initGin(s, server)
	preDeploy()
	ginEngine.GET("/", playgroundHandler())

	panic(ginEngine.Run(":" + port))
}
