package extensions

import (
	"context"
	"log/slog"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

func NewLogger(logger *slog.Logger) Logger {
	return Logger{logger: logger}
}

type Logger struct {
	logger *slog.Logger
}

func (l Logger) Validate(_ graphql.ExecutableSchema) error {
	return nil
}

func (l Logger) ExtensionName() string {
	return "Logger"
}

func (l Logger) InterceptField(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	if !l.logger.Handler().Enabled(ctx, slog.LevelDebug) {
		return next(ctx)
	}
	fc := graphql.GetFieldContext(ctx)
	ts := time.Now()
	resp, err := next(ctx)
	if fc.IsResolver {
		oc := graphql.GetOperationContext(ctx)
		l.logger.Debug(
			"resolver request",
			slog.String("raw-query", oc.RawQuery),
			slog.String("name", fc.Field.Name),
			slog.String("operation-name", oc.OperationName),
			slog.Float64("duration", time.Since(ts).Seconds()),
		)
	}
	return resp, err
}
