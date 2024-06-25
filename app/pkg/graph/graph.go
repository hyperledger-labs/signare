// Package graph contains the dependency management of the application.
package graph

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence/sql"
	"github.com/hyperledger-labs/signare/app/pkg/commons/validators"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
	"github.com/hyperledger-labs/signare/app/pkg/infra/rpcinfra"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/admin"

	"github.com/asaskevich/govalidator"

	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authentication/contextdefinition"
)

type GraphShared struct {
	useCasesGraph
}

type ApplicationGraph struct {
	config Config

	useCasesGraph       *useCasesGraph
	repositoriesGraph   *repositoriesGraph
	httpAPIGraph        *httpAPIGraph
	infraGraph          *infraGraph
	httpMiddlewareGraph *httpMiddlewareGraph
	rpcMiddlewareGraph  *rpcMiddlewareGraph
	librariesGraph      *librariesGraph
	metricRecorder      metricrecorder.MetricRecorder

	Shared GraphShared
}

func New(config Config) (*ApplicationGraph, error) {
	valid, err := govalidator.ValidateStruct(config)
	if err != nil || !valid {
		return nil, err
	}

	return &ApplicationGraph{
		config: config,
	}, nil
}

// Build builds the application graph
func (graph *ApplicationGraph) Build() {
	var err error

	// Validators
	validators.SetValidators()

	// Libraries
	graph.librariesGraph, err = initializeLibraries(graph.config)
	checkError(err)

	// Logger
	err = graph.librariesGraph.provideLogger(graph.config)
	checkError(err)

	// Repositories
	graph.repositoriesGraph, err = InitializeRepositories(graph.librariesGraph.persistenceFramework)
	checkError(err)

	// Metrics
	err = graph.initializeMetrics()
	checkError(err)

	// Use Cases
	graph.useCasesGraph, err = initializeUseCases(graph.repositoriesGraph, graph.metricRecorder, graph.config)
	checkError(err)

	// Infra
	graph.infraGraph, err = initializeInfra(graph.metricRecorder)
	checkError(err)

	// HTTP APIs
	authHeadersConfiguration := graph.getAuthHeadersConfiguration()
	graph.httpMiddlewareGraph, err = initializeHTTPMiddleware(graph.infraGraph, graph.useCasesGraph, graph.metricRecorder, authHeadersConfiguration)
	checkError(err)

	httpMiddleware := graph.httpMiddlewareGraph.HTTPMiddlewareFactory.Create()
	err = graph.infraGraph.mainHTTPRouter.RegisterMiddleware(httpMiddleware...)
	checkError(err)

	graph.rpcMiddlewareGraph, err = initializeRPCMiddleware(graph.infraGraph, graph.useCasesGraph, graph.metricRecorder, authHeadersConfiguration)
	checkError(err)

	rpcMiddleware := graph.rpcMiddlewareGraph.RPCMiddlewareFactory.Create()
	err = graph.infraGraph.rpcRouter.RegisterMiddleware(rpcMiddleware...)
	checkError(err)

	graph.httpAPIGraph, err = initializeHTTPAPI(graph.useCasesGraph, graph.infraGraph)
	checkError(err)

	// Metric Routes
	if graph.config.Libraries.Metrics != nil {
		err = graph.httpAPIGraph.publishMetricsHTTPListeners(graph.config.Libraries.Metrics.Prometheus, *graph.infraGraph)
		checkError(err)
	}
}

func (graph *ApplicationGraph) SetInitialSignerAdministrator(adminID string) (string, error) {
	listAdminsInput := admin.ListAdminsInput{}
	listAdminsOutput, err := graph.useCasesGraph.AdminUseCase.ListAdmins(context.Background(), listAdminsInput)
	if err != nil {
		return fmt.Sprintf("error when trying to create initial signer administrator '%s': %v", adminID, err), nil
	}
	for _, item := range listAdminsOutput.AdminCollection.Items {
		if item.ID == adminID {
			return fmt.Sprintf("initial signer administrator '%s' already existed -creation skipped", adminID), nil
		}
	}
	if len(listAdminsOutput.AdminCollection.Items) > 0 {
		if adminID == "" {
			return "signer administrators are already configured", nil
		}
		return "", errors.Internal().WithMessage("couldn't configure '%s' as initial signer administrator because there are already signer administrators configured", adminID)
	}

	description := "initial signer administrator configured from command line"
	createAdminInput := admin.CreateAdminInput{
		StandardID: entities.StandardID{
			ID: adminID,
		},
		Description: &description,
	}
	_, err = graph.useCasesGraph.AdminUseCase.CreateAdmin(context.Background(), createAdminInput)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return fmt.Sprintf("initial signer administrator '%s' already existed -creation skipped", adminID), nil
		}
		return "", errors.InternalFromErr(err).WithMessage("error creating initial signer administrator")
	}

	return fmt.Sprintf("initial signer administrator '%s' has been created", adminID), nil
}

// MainServer obtains the signare main HTTP server
func (graph *ApplicationGraph) MainServer() *httpinfra.DefaultHTTPRouter {
	if graph.infraGraph.mainHTTPRouter == nil {
		panic(fmt.Errorf("HTTP server not initialized"))
	}
	return graph.infraGraph.mainHTTPRouter
}

// RPCServer obtains the signare main RPC server
func (graph *ApplicationGraph) RPCServer() rpcinfra.RPCRouter {
	if graph.infraGraph.rpcRouter == nil {
		panic(fmt.Errorf("RPC server not initialized"))
	}
	return graph.infraGraph.rpcRouter
}

// MetricServer obtains the orchestrated metrics server. If the metrics server is nil and there is no error, it means that metrics are not used and therefore the metric server must not be deployed
func (graph *ApplicationGraph) MetricServer() *httpinfra.MetricsHTTPRouter {
	if graph.infraGraph.metricsHTTPRouter == nil {
		panic(fmt.Errorf("metrics HTTP server not initialized"))
	}
	return graph.infraGraph.metricsHTTPRouter
}

func (graph *ApplicationGraph) PersistenceFwConnection() sql.Connection {
	return graph.librariesGraph.persistenceConnection
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func (graph *ApplicationGraph) initializeMetrics() error {
	if graph.config.Libraries.Metrics != nil {
		metrics, err := initializePrometheusMetrics(graph.config)
		if err != nil {
			return err
		}
		graph.metricRecorder = metrics.metricRecorder
		return nil
	}
	dummyMetrics, err := InitializeDummyMetrics()
	if err != nil {
		return err
	}
	graph.metricRecorder = dummyMetrics.metricRecorder
	return nil
}

func (graph *ApplicationGraph) RPCMethods() []string {
	var methods []string
	if graph.infraGraph.rpcRouter == nil {
		return methods
	}
	return graph.infraGraph.rpcRouter.Methods()
}

func (graph *ApplicationGraph) UseCases() *GraphShared {
	return &GraphShared{
		useCasesGraph: *graph.useCasesGraph,
	}
}

func (graph *ApplicationGraph) getAuthHeadersConfiguration() contextdefinition.AuthHeadersConfiguration {
	var authHeadersConfiguration contextdefinition.AuthHeadersConfiguration
	if graph.config.RequestContextConfig != nil {
		authHeadersConfiguration.UserRequestHeader = graph.config.RequestContextConfig.UserHeaderKey
		authHeadersConfiguration.ApplicationRequestHeader = graph.config.RequestContextConfig.ApplicationHeaderKey
	}
	return authHeadersConfiguration
}
