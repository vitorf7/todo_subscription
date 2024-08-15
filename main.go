package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/vitorf7/todo_subscription/graph"
	"github.com/vitorf7/todo_subscription/graph/generated"
	"github.com/vitorf7/todo_subscription/internal/graph/extensions"
	"github.com/vitorf7/todo_subscription/internal/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "address",
		Usage:   "Service bind address",
		EnvVars: []string{"ADDRESS"},
		Value:   ":8888",
	},
}

const (
	appName        = "todo-graphQL"
	appDescription = "GraphQL for contact channels"
)

//go:generate go get -d github.com/99designs/gqlgen
//go:generate go run github.com/99designs/gqlgen

//nolint:funlen,cyclop,gocognit // main function.
func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Description = appDescription
	app.Flags = flags

	app.Action = func(cliCtx *cli.Context) error {
		logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}))
		slog.SetDefault(logger)

		srv := handler.New(generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					Logger: logger,
				},
			}))
		srv.AddTransport(transport.SSE{})
		// default server
		srv.AddTransport(transport.Options{})
		srv.AddTransport(transport.GET{})
		srv.AddTransport(transport.POST{})
		srv.AddTransport(transport.MultipartForm{})
		srv.SetQueryCache(lru.New(1000))
		srv.Use(extension.Introspection{})
		srv.Use(extension.AutomaticPersistedQuery{
			Cache: lru.New(100),
		})

		srv.Use(extensions.PrometheusInterceptor{})
		srv.Use(extensions.NewLogger(logger))
		srv.Use(extensions.Tracer{})

		router := mux.NewRouter()
		router.Use(
			middleware.Cors,
			middleware.GZIP,
			func(h http.Handler) http.Handler { return otelhttp.NewHandler(h, "graphQL") },
		)
		router.Handle("/graphql", srv)
		router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))

		logrus.Info("the server is listening on ", cliCtx.String("address"))
		server := &http.Server{
			Addr:              cliCtx.String("address"),
			Handler:           router,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       10 * time.Second,
		}
		return server.ListenAndServe()
	}

	if err := app.Run(os.Args); err != nil {
		logrus.WithError(err).Error("application failed")
	}
}
