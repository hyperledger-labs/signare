//go:build wireinject

package graph

import (
	"github.com/google/wire"
	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"
	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
	"github.com/hyperledger-labs/signare/app/pkg/infra/rpcinfra"
)

type infraGraph struct {
	httpAPIResponseHandler         httpinfra.HTTPResponseHandler
	defaultRPCInfraResponseHandler *rpcinfra.DefaultRPCInfraResponseHandler
	mainHTTPRouter                 *httpinfra.DefaultHTTPRouter
	metricsHTTPRouter              *httpinfra.MetricsHTTPRouter
	rpcRouter                      *rpcinfra.DefaultRPCRouter
}

var infraSet = wire.NewSet(
	wire.Struct(new(infraGraph), "*"),

	// Main HTTP HTTPRouter
	httpinfra.ProvideHTTPRouter,

	// Metrics HTTP HTTPRouter
	httpinfra.ProvideMetricsHTTPRouter,

	// HTTP metrics interface
	httpinfra.ProvideDefaultHTTPMetrics,
	wire.Bind(new(httpinfra.HTTPMetrics), new(*httpinfra.DefaultHTTPMetrics)),
	wire.Struct(new(httpinfra.DefaultHTTPMetricsOptions), "*"),

	// HTTP Response Handler
	httpinfra.ProvideDefaultHTTPResponseHandler,
	wire.Bind(new(httpinfra.HTTPResponseHandler), new(*httpinfra.DefaultHTTPResponseHandler)),
	wire.Struct(new(httpinfra.DefaultHTTPResponseHandlerOptions), "*"),

	// RPC Response Handler
	rpcinfra.ProvideDefaultRPCInfraResponseHandler,
	wire.Struct(new(rpcinfra.DefaultRPCInfraResponseHandlerOptions), "*"),

	// JSON-RPC RPCRouter
	rpcinfra.ProvideDefaultRPCRouter,
	wire.Struct(new(rpcinfra.DefaultRPCRouterOptions), "*"),
)

func initializeInfra(metricRecorder metricrecorder.MetricRecorder) (*infraGraph, error) {
	wire.Build(infraSet)
	return &infraGraph{}, nil
}
