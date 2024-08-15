package extensions

import (
	"context"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "graphql_requests_total",
		Help: "total graphql requests",
	}, []string{"code"})

	requestsHist = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "graphql_requests_seconds",
		Help: "time taken to resolve a gql query",
	})

	resolverTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "graphql_resolver_requests_total",
		Help: "total graphql resolver requests",
	}, []string{"name", "type", "code"})

	resolverHist = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "graphql_resolver_seconds",
		Help: "time taken to resolve a gql query",
	}, []string{"name", "type", "code"})
)

type PrometheusInterceptor struct{}

func (a PrometheusInterceptor) ExtensionName() string {
	return "Prometheus"
}

func (a PrometheusInterceptor) Validate(graphql.ExecutableSchema) error {
	return nil
}

func (a PrometheusInterceptor) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	resp := next(ctx)
	errs := graphql.GetErrors(ctx)
	switch {
	case len(errs) > 0:
		requestsTotal.WithLabelValues("error").Inc()
	default:
		requestsTotal.WithLabelValues("ok").Inc()
	}
	requestsHist.Observe(time.Since(graphql.GetOperationContext(ctx).Stats.OperationStart).Seconds())
	return resp
}

func (a PrometheusInterceptor) InterceptField(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	fc := graphql.GetFieldContext(ctx)
	resp, err := next(ctx)
	if fc.IsResolver {
		labels := make([]string, 0, 3)
		labels = append(labels, fc.Field.Name, fc.Object)
		oc := graphql.GetOperationContext(ctx)
		if err != nil || len(graphql.GetErrors(ctx)) > 0 {
			labels = append(labels, "error")
			resolverTotal.WithLabelValues(labels...).Inc()
			resolverHist.WithLabelValues(labels...).Observe(time.Since(oc.Stats.OperationStart).Seconds())
		} else {
			labels = append(labels, "ok")
			resolverTotal.WithLabelValues(labels...).Inc()
			resolverHist.WithLabelValues(labels...).Observe(time.Since(oc.Stats.OperationStart).Seconds())
		}
	}
	return resp, err
}
