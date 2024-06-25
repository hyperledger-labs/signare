package httpinfra

// MetricsHTTPRouter HTTP router used for metrics
type MetricsHTTPRouter struct {
	// DefaultHTTPRouter holds the router to handle HTTP incoming connections.
	DefaultHTTPRouter
}

// ProvideMetricsHTTPRouter creates a MetricsHTTPRouter
func ProvideMetricsHTTPRouter() *MetricsHTTPRouter {
	r := ProvideHTTPRouter()
	return &MetricsHTTPRouter{
		DefaultHTTPRouter: *r,
	}
}
