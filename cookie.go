package torii

type (
	cookie struct {
		name  string
		value string
	}
	cookieOption func(c *cookie)
)

func (c *cookie) String() string {
	return ""
}

// SetCookie sets cookie to the response.
func SetCookie(response Response, name, value string, options ...cookieOption) Response {
	c := cookie{
		name:  name,
		value: value,
	}
	for _, option := range options {
		option(&c)
	}
	return Response{
		Status: response.Status,
		Head: append(
			response.Head,
			Header{
				Key:   "Set-Cookie",
				Value: c.String(),
			},
		),
		Body: response.Body,
	}
}
