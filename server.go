package spring

import (
	"fmt"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
)

func (s *Spring) RunServer(defaultPort uint, server graphql.ExecutableSchema, ginInitHandler GinInitHandler) {
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", defaultPort)
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: InitGin(server, ginInitHandler),
	}
	func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
		s.done <- true
	}()
	s.await()
}
