package torii

import (
	"bufio"
	"log"
	"net/http"
)

type handler struct {
	gate Gate

	logger     *log.Logger
	bufferSize uint
}

type handlerOption func(h *handler)

// Handler returns Gate wrapper to suit http.Handler.
func Handler(gate Gate, options ...handlerOption) http.Handler {
	if gate == nil {
		panic("no gate provided to Handler")
	}
	handler := &handler{
		gate:       gate,
		bufferSize: 2048,
	}
	for _, applyOptionTo := range options {
		applyOptionTo(handler)
	}
	return handler
}

// ServeHTTP implements http.Handler.
func (h *handler) ServeHTTP(out http.ResponseWriter, in *http.Request) {
	response := h.gate.Respond(
		in.Context(),
		Request{
			Method: (Method)(in.Method),
			Target: in.URL,
			Headers: (func() []Header {
				var headers []Header
				for key, values := range in.Header {
					for _, value := range values {
						headers = append(
							headers,
							Header{
								Key:   key,
								Value: value,
							},
						)
					}
				}
				return headers
			})(),
			Body: in.Body,
		},
	)
	if headers := response.Headers; headers != nil {
		for _, header := range headers {
			out.Header().Add(header.Key, header.Value)
		}
	}
	out.WriteHeader((int)(response.Status))
	if body := response.Body; body != nil {
		defer (func() {
			if err := response.Body.Close(); err != nil {
				if logger := h.logger; logger != nil {
					logger.Printf("Failed to close response body: %s", err.Error())
				}
			}
		})()
		if _, err := bufio.NewReaderSize(body, (int)(h.bufferSize)).WriteTo(out); err != nil {
			if logger := h.logger; logger != nil {
				logger.Printf("Failed to write response body: %s", err.Error())
			}
		}
	} else {
		if _, err := out.Write(([]byte)(http.StatusText((int)(response.Status)))); err != nil {
			if logger := h.logger; logger != nil {
				logger.Printf("Failed to write response body: %s", err.Error())
			}
		}
	}
}

// HandlerOptionWithLogger sets the logger to the handler.
func HandlerOptionWithLogger(logger *log.Logger) handlerOption {
	if logger == nil {
		panic("logger must not be nil (if you don't need a logger, just omit this option)")
	}
	return (handlerOption)(func(h *handler) {
		h.logger = logger
	})
}

// HandlerOptionWithBufferSize sets buffer size for the handler.
//
// The buffer is used to write response bodies.
//
// This value defaults to 2048.
func HandlerOptionWithBufferSize(value uint) handlerOption {
	return (handlerOption)(func(h *handler) {
		h.bufferSize = value
	})
}
