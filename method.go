package torii

import "net/http"

// Method is http method.
type Method string

// Type-safe HTTP methods.
const (
	MethodConnect = (Method)(http.MethodConnect)
	MethodDelete  = (Method)(http.MethodDelete)
	MethodHead    = (Method)(http.MethodHead)
	MethodGet     = (Method)(http.MethodGet)
	MethodPost    = (Method)(http.MethodPost)
	MethodPut     = (Method)(http.MethodPut)
	MethodPatch   = (Method)(http.MethodPatch)
	MethodOptions = (Method)(http.MethodOptions)
	MethodTrace   = (Method)(http.MethodTrace)
)

// String returns string representation of the Method.
func (m Method) String() string {
	return (string)(m)
}
