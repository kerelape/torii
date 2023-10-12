package torii

import (
	"context"
)

// GateNotFound is a gate that always responds with status NotFound.
var GateNotFound Gate = Func(func(context.Context, Request) Response {
	return Response{
		Status: StatusNotFound,
	}
})
