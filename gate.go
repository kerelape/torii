package torii

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

// Gate represents an endpoint.
type Gate interface {

	// Respond returns response for the request.
	Respond(context.Context, Request) Response
}

type (
	// Header of a message.
	Header http.Header

	// Response message.
	Response struct {
		Status Status
		Header Header
		Body   io.ReadCloser
	}

	// Request message.
	Request struct {
		Target *url.URL
		Method Method
		Header Header
		Body   io.Reader
	}
)
