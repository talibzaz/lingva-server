package server

import (
	"context"
	"net/http"
	"lingva/gql"

	"github.com/rs/cors"
	"github.com/vektah/gqlgen/graphql"
	"github.com/vektah/gqlgen/handler"
	"github.com/vektah/gqlparser/gqlerror"
	"github.com/gorilla/mux"
	"lingva/auth"
)

// GraphQL Server which implements the ResolverRoot interface present inside generated.go.
type server struct {

}

// NewHTTPServer - returns a new *http.Server instance bootstrapped with our own GraphQL server.
func NewHTTPServer(port string) (*http.Server, error) {
	s := &server{}

	mux := mux.NewRouter()

	mux.Handle("/", handler.Playground("LINGVA Playground", "/graphql"))
	config := gql.Config{Resolvers: s}

	mux.Handle("/graphql", auth.ValidateMiddleware(handler.GraphQL(
		gql.NewExecutableSchema(config),
		handler.ErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
			return &gqlerror.Error{
				Message: err.Error(),
				Path:    graphql.GetResolverContext(ctx).Path(),
			}
		}),
	)))

	mux.HandleFunc("/login", Login).Methods("POST")

	mux.HandleFunc("/admin-login", AdminLogin).Methods("POST")

	mux.HandleFunc("/image/{imageName}", ImageHandler).Methods("GET")

	return &http.Server{
		Addr:    ":" + port,
		Handler: cors.AllowAll().Handler(mux),
	}, nil
}