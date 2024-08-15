package extensions

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct{}

func (Tracer) ExtensionName() string {
	return "Tracer"
}

func (t Tracer) Validate(_ graphql.ExecutableSchema) error {
	return nil
}

func (t Tracer) InterceptField(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	fc := graphql.GetFieldContext(ctx)
	if fc.IsResolver {
		var span trace.Span
		ctx, span = otel.Tracer("").Start(ctx, "Resolve",
			trace.WithAttributes(
				attribute.String("graphql.resolver", fc.Field.Name),
				attribute.String("graphql.type", fc.Object),
			))
		defer span.End()
	}

	return next(ctx)
}
