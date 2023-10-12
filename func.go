package torii

import "context"

// Func is Gate func implementation.
type Func func(context.Context, Request) Response

// Respond implements Gate.
func (f Func) Respond(ctx context.Context, request Request) Response {
	return f(ctx, request)
}
