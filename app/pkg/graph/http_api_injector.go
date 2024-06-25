//go:build wireinject

package graph

import (
	"github.com/google/wire"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/httpin"
	"github.com/hyperledger-labs/signare/app/pkg/adapters/rpcin"
	generatedhttpinfra "github.com/hyperledger-labs/signare/app/pkg/infra/generated/httpinfra"
	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
	"github.com/hyperledger-labs/signare/app/pkg/infra/rpcinfra"

	"github.com/hyperledger-labs/signare/app/pkg/infra/metricshttpinfra"
)

type httpAPIGraph struct {
	// REST
	adminAPIRoutesPublished       generatedhttpinfra.AdminAPIRoutesPublished
	applicationAPIRoutesPublished generatedhttpinfra.ApplicationAPIRoutesPublished

	// JSON-RPC
	rpcAPIRoutesPublished rpcinfra.JSONRPCAPIRoutesPublished
}

var httpAPISet = wire.NewSet(
	wire.Struct(new(httpAPIGraph), "*"),

	/**************/
	/*    REST   */
	/*************/

	// Main Router
	provideMainRouter,

	// Admin API Adapter
	httpin.ProvideDefaultAdminAPIAdapter,
	wire.Bind(new(generatedhttpinfra.AdminAPIAdapter), new(*httpin.DefaultAdminAPIAdapter)),
	wire.Struct(new(httpin.DefaultAdminAPIAdapterOptions), "*"),

	// Admin API Handler
	generatedhttpinfra.NewDefaultAdminAPIHTTPHandler,
	wire.Bind(new(generatedhttpinfra.AdminAPIHTTPHandler), new(*generatedhttpinfra.DefaultAdminAPIHTTPHandler)),
	wire.Struct(new(generatedhttpinfra.DefaultAdminAPIHTTPHandlerOptions), "*"),

	// Admin API Routes
	generatedhttpinfra.ProvideAdminAPIRoutes,
	wire.Bind(new(httpinfra.HTTPRouter), new(*httpinfra.DefaultHTTPRouter)),
	wire.Struct(new(generatedhttpinfra.AdminAPIPublisherOptions), "*"),

	// Application API Adapter
	httpin.ProvideDefaultApplicationAPIAdapter,
	wire.Bind(new(generatedhttpinfra.ApplicationAPIAdapter), new(*httpin.DefaultApplicationAPIAdapter)),
	wire.Struct(new(httpin.DefaultApplicationAPIAdapterOptions), "*"),

	// Application API Handler
	generatedhttpinfra.NewDefaultApplicationAPIHTTPHandler,
	wire.Bind(new(generatedhttpinfra.ApplicationAPIHTTPHandler), new(*generatedhttpinfra.DefaultApplicationAPIHTTPHandler)),
	wire.Struct(new(generatedhttpinfra.DefaultApplicationAPIHTTPHandlerOptions), "*"),

	// Application API Routes
	generatedhttpinfra.ProvideApplicationAPIRoutes,
	wire.Struct(new(generatedhttpinfra.ApplicationAPIPublisherOptions), "*"),

	/*****************/
	/*   JSON-RPC  */
	/***************/

	// JSON-RPC API Adapter
	rpcin.NewDefaultAPIAdapter,
	wire.Bind(new(rpcinfra.JSONRPCAPIAdapter), new(*rpcin.DefaultAPIAdapter)),
	wire.Struct(new(rpcin.DefaultAPIAdapterOptions), "*"),

	// JSON-RPC API Handler
	rpcinfra.NewDefaultJSONRPCAPIHandler,
	wire.Bind(new(rpcinfra.JSONRPCAPIHandler), new(*rpcinfra.DefaultJSONRPCAPIHandler)),
	wire.Struct(new(rpcinfra.DefaultJSONRPCAPIHandlerOptions), "*"),

	// JSON-RPC API Methods
	rpcinfra.ProvideJSONRPCMethods,
	wire.Struct(new(rpcinfra.JSONRPCAPIPublisherOptions), "*"),
)

func initializeHTTPAPI(
	useCases *useCasesGraph,
	infra *infraGraph,
) (*httpAPIGraph, error) {
	wire.Build(httpAPISet,
		wire.FieldsOf(new(*useCasesGraph),
			"ApplicationUseCase",
			"AccountUseCase",
			"UserUseCase",
			"AdminUseCase",
			"HSMModuleUseCase",
			"HSMSlotUseCase",
			"HSMConnector",
			"HSMConnectionResolver",
		),
		wire.FieldsOf(new(*infraGraph),
			"httpAPIResponseHandler",
			"rpcRouter",
		),
		wire.Bind(new(rpcinfra.RPCRouter), new(*rpcinfra.DefaultRPCRouter)),
	)

	return &httpAPIGraph{}, nil
}

func provideMainRouter(infra *infraGraph) *httpinfra.DefaultHTTPRouter {
	return infra.mainHTTPRouter
}

func (g *httpAPIGraph) publishMetricsHTTPListeners(config PrometheusConfig, infra infraGraph) error {
	options := metricshttpinfra.PrometheusMetricsHTTPOptions{
		HTTPInfra: infra.metricsHTTPRouter,
		PrometheusConfig: metricshttpinfra.PrometheusConfig{
			Path:                config.Path,
			MaxRequestsInFlight: config.MaxRequestsInFlight,
			TimeoutInMillis:     config.TimeoutInMillis,
		},
	}
	_, err := metricshttpinfra.ProvideMetricsHTTP(options)
	if err != nil {
		return err
	}
	return nil
}
