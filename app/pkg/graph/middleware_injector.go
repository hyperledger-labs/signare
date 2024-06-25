//go:build wireinject

package graph

import (
	"github.com/google/wire"
	embedded "github.com/hyperledger-labs/signare/app"
	"github.com/hyperledger-labs/signare/app/pkg/adapters/httpmiddlewarein/pepin"
	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/infile/pipinfile"
	"github.com/hyperledger-labs/signare/app/pkg/adapters/usecaseadapters/pip"
	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"
	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authentication"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authentication/contextdefinition"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authentication/contextdefinition/httpcontextdefinition"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authentication/contextdefinition/rpccontextdefinition"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authentication/contextvalidation"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authorization"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authorization/pep"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/entrypoint/rpcbatchrequestsupport"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/telemetry"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/telemetry/tracer"
	"github.com/hyperledger-labs/signare/app/pkg/infra/rpcinfra"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/authorization/pdp"
)

func ProvidePolicyInformationPointYAMLOutputAdapter() *pipinfile.DefaultRBACActionsPolicyInformationPointYAMLOutputAdapter {
	policyDecisionPointOutputAdapterOptions := pipinfile.DefaultRBACActionsPolicyInformationPointYAMLOutputAdapterOptions{
		FileSystem: embedded.RBACFiles,
		BasePath:   "include/rbac",
	}
	policyInformationPointOutputAdapter, err := pipinfile.ProvideDefaultRBACActionsPolicyInformationPointYAMLOutputAdapter(policyDecisionPointOutputAdapterOptions)
	checkError(err)

	return policyInformationPointOutputAdapter
}

type httpMiddlewareGraph struct {
	HTTPMiddlewareFactory *middleware.HTTPMiddlewareFactory
}

var httpMiddlewareSet = wire.NewSet(
	wire.Struct(new(httpMiddlewareGraph), "*"),

	httpcontextdefinition.ProvideHTTPContextDefinition,
	wire.Bind(new(contextdefinition.ContextDefinition), new(*httpcontextdefinition.HTTPContextDefinition)),
	wire.Struct(new(httpcontextdefinition.HTTPContextDefinitionOptions), "*"),

	contextvalidation.ProvideRequestContextValidation,
	wire.Struct(new(contextvalidation.RequestContextValidationOptions), "*"),

	pip.ProvideDefaultAccountsPIPAdapter,
	wire.Bind(new(pdp.AccountsPolicyInformationPort), new(*pip.DefaultAccountsPIPAdapter)),
	wire.Struct(new(pip.DefaultAccountsPIPAdapterOptions), "*"),

	pip.ProvideDefaultAdminsPIPAdapter,
	wire.Bind(new(pdp.AdminsPolicyInformationPort), new(*pip.DefaultAdminsPIPAdapter)),
	wire.Struct(new(pip.DefaultAdminsPIPAdapterOptions), "*"),

	ProvidePolicyInformationPointYAMLOutputAdapter,
	wire.Bind(new(pdp.ActionsPolicyInformationPointPort), new(*pipinfile.DefaultRBACActionsPolicyInformationPointYAMLOutputAdapter)),

	pip.ProvideDefaultUsersPIPAdapter,
	wire.Bind(new(pdp.UsersPolicyInformationPort), new(*pip.DefaultUsersPIPAdapter)),
	wire.Struct(new(pip.DefaultUsersPIPAdapterOptions), "*"),

	pdp.ProvideDefaultPolicyDecisionPointUseCase,
	wire.Bind(new(pdp.PolicyDecisionPointUseCase), new(*pdp.DefaultPolicyDecisionPointUseCase)),
	wire.Struct(new(pdp.DefaultPolicyDecisionPointUseCaseOptions), "*"),

	pepin.ProvideUserPolicyDecisionPointAdapter,
	wire.Bind(new(pep.UserPolicyDecisionPointPort), new(*pepin.DefaultUserPolicyDecisionPointAdapter)),
	wire.Struct(new(pepin.DefaultUserPolicyDecisionPointAdapterOptions), "*"),

	pepin.ProvideDefaultAccountUserPolicyDecisionPointAdapter,
	wire.Bind(new(pep.AccountUserPolicyDecisionPointPort), new(*pepin.DefaultAccountUserPolicyDecisionPointAdapter)),
	wire.Struct(new(pepin.DefaultAccountUserPolicyDecisionPointAdapterOptions), "*"),

	pep.ProvideHTTPPolicyEnforcementPoint,
	wire.Struct(new(pep.HTTPPolicyEnforcementPointOptions), "*"),

	pep.ProvideRPCPolicyEnforcementPoint,
	wire.Struct(new(pep.RPCPolicyEnforcementPointOptions), "*"),

	authorization.ProvideAuthorizationMiddleware,
	wire.Struct(new(authorization.AuthorizationMiddlewareOptions), "*"),

	authentication.ProvideAuthenticationMiddleware,
	wire.Struct(new(authentication.AuthenticationMiddlewareOptions), "*"),

	middleware.ProvideHTTPMiddlewareFactory,
	wire.Struct(new(middleware.HTTPMiddlewareFactoryOptions), "*"),

	telemetry.ProvideTelemetryMiddleware,
	wire.Struct(new(telemetry.TelemetryMiddlewareOptions), "*"),

	tracer.ProvideHTTPContextTracer,
	wire.Struct(new(tracer.HTTPContextTracerOptions), "*"),
)

func initializeHTTPMiddleware(
	infra *infraGraph,
	useCases *useCasesGraph,
	metricRecorder metricrecorder.MetricRecorder,
	configuration contextdefinition.AuthHeadersConfiguration,
) (*httpMiddlewareGraph, error) {
	wire.Build(httpMiddlewareSet,
		wire.FieldsOf(new(*infraGraph),
			"httpAPIResponseHandler",
			"mainHTTPRouter",
		),
		wire.FieldsOf(new(*useCasesGraph),
			"UserUseCase",
			"AccountUseCase",
			"AdminUseCase",
		),
		wire.Bind(new(httpinfra.HTTPRouter), new(*httpinfra.DefaultHTTPRouter)),
	)

	return &httpMiddlewareGraph{}, nil
}

type rpcMiddlewareGraph struct {
	RPCMiddlewareFactory *middleware.RPCMiddlewareFactory
}

var rpcMiddlewareSet = wire.NewSet(

	wire.Struct(new(rpcMiddlewareGraph), "*"),

	contextvalidation.ProvideRequestContextValidation,
	wire.Struct(new(contextvalidation.RequestContextValidationOptions), "*"),

	rpccontextdefinition.ProvideRPCContextDefinitionFromHeaders,
	wire.Bind(new(contextdefinition.ContextDefinition), new(*rpccontextdefinition.RPCContextDefinition)),
	wire.Struct(new(rpccontextdefinition.RPCContextDefinitionOptions), "*"),

	pip.ProvideDefaultAccountsPIPAdapter,
	wire.Bind(new(pdp.AccountsPolicyInformationPort), new(*pip.DefaultAccountsPIPAdapter)),
	wire.Struct(new(pip.DefaultAccountsPIPAdapterOptions), "*"),

	pip.ProvideDefaultAdminsPIPAdapter,
	wire.Bind(new(pdp.AdminsPolicyInformationPort), new(*pip.DefaultAdminsPIPAdapter)),
	wire.Struct(new(pip.DefaultAdminsPIPAdapterOptions), "*"),

	ProvidePolicyInformationPointYAMLOutputAdapter,
	wire.Bind(new(pdp.ActionsPolicyInformationPointPort), new(*pipinfile.DefaultRBACActionsPolicyInformationPointYAMLOutputAdapter)),

	pip.ProvideDefaultUsersPIPAdapter,
	wire.Bind(new(pdp.UsersPolicyInformationPort), new(*pip.DefaultUsersPIPAdapter)),
	wire.Struct(new(pip.DefaultUsersPIPAdapterOptions), "*"),

	pdp.ProvideDefaultPolicyDecisionPointUseCase,
	wire.Bind(new(pdp.PolicyDecisionPointUseCase), new(*pdp.DefaultPolicyDecisionPointUseCase)),
	wire.Struct(new(pdp.DefaultPolicyDecisionPointUseCaseOptions), "*"),

	pepin.ProvideUserPolicyDecisionPointAdapter,
	wire.Bind(new(pep.UserPolicyDecisionPointPort), new(*pepin.DefaultUserPolicyDecisionPointAdapter)),
	wire.Struct(new(pepin.DefaultUserPolicyDecisionPointAdapterOptions), "*"),

	pepin.ProvideDefaultAccountUserPolicyDecisionPointAdapter,
	wire.Bind(new(pep.AccountUserPolicyDecisionPointPort), new(*pepin.DefaultAccountUserPolicyDecisionPointAdapter)),
	wire.Struct(new(pepin.DefaultAccountUserPolicyDecisionPointAdapterOptions), "*"),

	pep.ProvideHTTPPolicyEnforcementPoint,
	wire.Struct(new(pep.HTTPPolicyEnforcementPointOptions), "*"),

	pep.ProvideRPCPolicyEnforcementPoint,
	wire.Struct(new(pep.RPCPolicyEnforcementPointOptions), "*"),

	authorization.ProvideAuthorizationMiddleware,
	wire.Struct(new(authorization.AuthorizationMiddlewareOptions), "*"),

	authentication.ProvideAuthenticationMiddleware,
	wire.Struct(new(authentication.AuthenticationMiddlewareOptions), "*"),

	rpcbatchrequestsupport.ProvideRPCBatchRequestSupportMiddleware,
	wire.Struct(new(rpcbatchrequestsupport.RPCBatchRequestSupportMiddlewareOptions), "*"),

	middleware.ProvideRPCMiddlewareFactory,
	wire.Struct(new(middleware.RPCMiddlewareFactoryOptions), "*"),

	telemetry.ProvideTelemetryMiddleware,
	wire.Struct(new(telemetry.TelemetryMiddlewareOptions), "*"),

	tracer.ProvideHTTPContextTracer,
	wire.Struct(new(tracer.HTTPContextTracerOptions), "*"),
)

func initializeRPCMiddleware(
	infra *infraGraph,
	useCases *useCasesGraph,
	metricRecorder metricrecorder.MetricRecorder,
	configuration contextdefinition.AuthHeadersConfiguration,
) (*rpcMiddlewareGraph, error) {
	wire.Build(rpcMiddlewareSet,
		wire.FieldsOf(new(*infraGraph),
			"defaultRPCInfraResponseHandler",
			"rpcRouter",
		),
		wire.FieldsOf(new(*useCasesGraph),
			"UserUseCase",
			"AccountUseCase",
			"AdminUseCase",
		),
		wire.Bind(new(rpcinfra.RPCRouter), new(*rpcinfra.DefaultRPCRouter)),
		wire.Bind(new(httpinfra.HTTPResponseHandler), new(*rpcinfra.DefaultRPCInfraResponseHandler)),
	)

	return &rpcMiddlewareGraph{}, nil
}
