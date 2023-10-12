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
	Header struct {
		Key   string
		Value string
	}

	// Response message.
	Response struct {
		Status  Status
		Headers []Header
		Cookies []http.Cookie
		Body    io.ReadCloser
	}

	// Request message.
	Request struct {
		Target  *url.URL
		Method  Method
		Headers []Header
		Cookies []http.Cookie
		Body    io.Reader
	}
)
