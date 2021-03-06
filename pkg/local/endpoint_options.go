package local

import (
	"context"
)

// endpointOptions controls the override behavior for a local endpoint.
type endpointOptions struct {
	cachePathCallback   func(string, bool) (string, error)
	stagingRootCallback func(string, bool) (string, error)
	watchingMechanism   func(context.Context, string, chan<- struct{})
}

// EndpointOption is the interface for specifying endpoint options. It cannot be
// constructed or implemented directly, only by one of option constructors
// provided by this package.
type EndpointOption interface {
	// apply modifies the provided endpoint options configuration in accordance
	// with the option.
	apply(*endpointOptions)
}

// functionEndpointOption is an implementation of EndpointOption that adapts a
// simple closure function to be an EndpointOption.
type functionEndpointOption struct {
	applier func(*endpointOptions)
}

// newFunctionEndpointOption creates a new EndpointOption using a simple closure
// function.
func newFunctionEndpointOption(applier func(*endpointOptions)) EndpointOption {
	return &functionEndpointOption{applier}
}

// apply implements EndpointOption.apply for functionEndpointOption.
func (o *functionEndpointOption) apply(options *endpointOptions) {
	o.applier(options)
}

// WithCachePathCallback overrides the function that the endpoint uses to
// compute cache storage paths. The specified callback will be provided with two
// arguments: the session identifier (a UUID) and a boolean indicating whether
// or not this is the alpha endpoint (if false, it's the beta endpoint). The
// function should return a path that is consistent but unique in terms of these
// two arguments.
func WithCachePathCallback(callback func(string, bool) (string, error)) EndpointOption {
	return newFunctionEndpointOption(func(options *endpointOptions) {
		options.cachePathCallback = callback
	})
}

// WithStagingRootCallback overrides the function that the endpoint uses to
// compute staging root paths. The specified callback will be provided with two
// arguments: the session identifier (a UUID) and a boolean indicating whether
// or not this is the alpha endpoint (if false, it's the beta endpoint). The
// function should return a path that is consistent but unique in terms of these
// two arguments. The path may exist, but if it does must be a directory.
func WithStagingRootCallback(callback func(string, bool) (string, error)) EndpointOption {
	return newFunctionEndpointOption(func(options *endpointOptions) {
		options.stagingRootCallback = callback
	})
}

// WithWatchingMechanism overrides the filesystem watching function that the
// endpoint uses to monitor for filesystem changes. The specified function will
// be provided with three arguments: a context to indicate watch cancellation,
// the path to be watched (recursively), and an events channel that should be
// populated in a non-blocking fashion every time an event occurs. If an error
// occurs during watching, the event channel should be closed. It should also be
// closed on cancellation.
func WithWatchingMechanism(callback func(context.Context, string, chan<- struct{})) EndpointOption {
	return newFunctionEndpointOption(func(options *endpointOptions) {
		options.watchingMechanism = callback
	})
}
