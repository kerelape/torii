package torii

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

type (
	fork struct {
		routes   map[string]map[Method]Gate
		fallback Gate
	}
	endpoint func() (string, Method, Gate)
)

// Endpoint returns a fork endpoint.
//
// The name must start wtih a slash (/) and not contain:
//   - new-line characters (\n);
//   - reset characters (\r);
//   - whitespace characters;
//   - slashes (excpet leading).
func Endpoint(method Method, name string, gate Gate) endpoint {
	if name == "" {
		panic("endpoint name cannot be empty")
	}
	if strings.ContainsRune(name, '\n') {
		panic("endpoint name must not contain new-line symbols (\\n)")
	}
	if strings.ContainsRune(name, '\r') {
		panic("endpoint name must not contain reset symbols (\\r)")
	}
	if strings.ContainsRune(name, ' ') {
		panic("endpoint name must not contain whitespace sybmols")
	}
	if !strings.HasPrefix(name, "/") {
		panic("endpoint name must start with a slash (/)")
	}
	if strings.ContainsRune(strings.TrimPrefix(name, "/"), '/') {
		panic("endpoint name must not contain slashes (/) expect for prefix")
	}
	return func() (string, Method, Gate) {
		return strings.Trim(name, "/"), method, gate
	}
}

// Fork creates a fork gate and returns it.
func Fork(fallback Gate, endpoints ...endpoint) Gate {
	if fallback == nil {
		panic("fallback cannot be nil")
	}
	f := fork{
		routes:   make(map[string]map[Method]Gate, len(endpoints)),
		fallback: fallback,
	}
	for _, endpoint := range endpoints {
		name, method, gate := endpoint()
		if _, ok := f.routes[name]; !ok {
			f.routes[name] = make(map[Method]Gate)
		}
		if _, ok := f.routes[name][method]; ok {
			panic(fmt.Sprintf("%s /%s already exists", method, name))
		}
		f.routes[name][method] = gate
	}
	return f
}

// Respond implements Gate.
func (f fork) Respond(ctx context.Context, request Request) Response {
	segments := strings.Split(strings.Trim(request.Target.Path, "/"), "/")
	if route, ok := f.routes[segments[0]]; ok {
		if gate, ok := route[request.Method]; ok {
			target, targetError := url.ParseRequestURI("/" + strings.Join(segments[1:], "/"))
			if targetError != nil {
				panic(targetError)
			}
			return gate.Respond(
				ctx,
				Request{
					Target: target,
					Method: request.Method,
					Header: request.Header,
					Body:   request.Body,
				},
			)
		} else {
			return Response{
				Status: StatusMethodNotAllowed,
			}
		}
	}
	return f.fallback.Respond(ctx, request)
}
