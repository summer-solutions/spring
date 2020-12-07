package spring

import (
	"fmt"
	"os"

	"github.com/99designs/gqlgen/graphql"
)

func (s *Spring) RunServer(defaultPort uint, server graphql.ExecutableSchema, ginInitHandler GinInitHandler) {
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", defaultPort)
	}
	ginEngine := initGin(server, ginInitHandler)
	panic(ginEngine.Run(":" + port))
}
