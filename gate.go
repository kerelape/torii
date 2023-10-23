package torii

import (
	"context"
	"io"
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

	// Head is a collection of headers.
	Head []Header

	// Response message.
	Response struct {
		Status Status
		Head   Head
		Body   io.ReadCloser
	}

	// Request message.
	Request struct {
		Target *url.URL
		Method Method
		Head   Head
		Body   io.Reader
	}
)

// Values returns all values set by the key.
//
// This method always returns a non-nil slice.
func (h Head) Values(key string) []string {
	var values []string
	for _, header := range h {
		if header.Key == key {
			values = append(values, header.Value)
		}
	}
	return values
}
