// Code generated by mga tool. DO NOT EDIT.
package clusterdriver

import (
	"github.com/banzaicloud/pipeline/internal/cluster"
	"github.com/go-kit/kit/endpoint"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	kitxendpoint "github.com/sagikazarmark/kitx/endpoint"
)

// Endpoint name constants
const (
	CreateNodePoolEndpoint = "cluster.CreateNodePool"
	DeleteClusterEndpoint  = "cluster.DeleteCluster"
	DeleteNodePoolEndpoint = "cluster.DeleteNodePool"
)

// Endpoints collects all of the endpoints that compose the underlying service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	CreateNodePool endpoint.Endpoint
	DeleteCluster  endpoint.Endpoint
	DeleteNodePool endpoint.Endpoint
}

// MakeEndpoints returns a(n) Endpoints struct where each endpoint invokes
// the corresponding method on the provided service.
func MakeEndpoints(service cluster.Service, middleware ...endpoint.Middleware) Endpoints {
	mw := kitxendpoint.Combine(middleware...)

	return Endpoints{
		CreateNodePool: kitxendpoint.OperationNameMiddleware(CreateNodePoolEndpoint)(mw(MakeCreateNodePoolEndpoint(service))),
		DeleteCluster:  kitxendpoint.OperationNameMiddleware(DeleteClusterEndpoint)(mw(MakeDeleteClusterEndpoint(service))),
		DeleteNodePool: kitxendpoint.OperationNameMiddleware(DeleteNodePoolEndpoint)(mw(MakeDeleteNodePoolEndpoint(service))),
	}
}

// TraceEndpoints returns a(n) Endpoints struct where each endpoint is wrapped with a tracing middleware.
func TraceEndpoints(endpoints Endpoints) Endpoints {
	return Endpoints{
		CreateNodePool: kitoc.TraceEndpoint("cluster.CreateNodePool")(endpoints.CreateNodePool),
		DeleteCluster:  kitoc.TraceEndpoint("cluster.DeleteCluster")(endpoints.DeleteCluster),
		DeleteNodePool: kitoc.TraceEndpoint("cluster.DeleteNodePool")(endpoints.DeleteNodePool),
	}
}