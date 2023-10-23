package cookie

import (
	"strconv"
	"strings"
	"time"

	"github.com/kerelape/torii"
	"github.com/kerelape/torii/internal/option"
)

// SameSitePolicy is cookie same-site policy.
type SameSitePolicy string

const (
	SameSitePolicyStrict = (SameSitePolicy)("Strict")
	SameSitePolicyLax    = (SameSitePolicy)("Lax")
	SameSitePolicyNone   = (SameSitePolicy)("None")
)

type (
	cookie struct {
		name  string
		value string

		httpOnly bool
		secure   bool
		expires  option.Option[time.Time]
		maxAge   option.Option[time.Duration]
		path     string
		sameSite SameSitePolicy
		domain   string
	}
	cookieOption func(c *cookie)
)

func (c *cookie) String() string {
	elems := make([]string, 0, 8)
	elems = append(elems, c.name+"="+c.value)
	if c.httpOnly {
		elems = append(elems, "HttpOnly")
	}
	if c.secure {
		elems = append(elems, "Secure")
	}
	if c.expires.HasValue() {
		elems = append(elems, "Expires="+c.expires.Value().Format(time.RFC1123))
	}
	if c.maxAge.HasValue() {
		elems = append(elems, "Max-Age="+strconv.Itoa((int)(c.maxAge.Value().Seconds())))
	}
	if c.domain != "" {
		elems = append(elems, "Domain="+c.domain)
	}
	if c.path != "" {
		elems = append(elems, "Path="+c.path)
	}
	if c.sameSite != "" {
		elems = append(elems, "SameSite="+(string)(c.sameSite))
	}
	return strings.Join(elems, "; ")
}

// Set sets cookie to the response.
func Set(response torii.Response, name, value string, options ...cookieOption) torii.Response {
	c := cookie{
		name:  name,
		value: value,
	}
	for _, option := range options {
		option(&c)
	}
	return torii.Response{
		Status: response.Status,
		Head: append(
			response.Head,
			torii.Header{
				Key:   "Set-Cookie",
				Value: c.String(),
			},
		),
		Body: response.Body,
	}
}

// HttpOnly makes the cookie http-only.
func HttpOnly() cookieOption {
	return func(c *cookie) {
		c.httpOnly = true
	}
}

// Secure makes the cookie secure.
func Secure() cookieOption {
	return func(c *cookie) {
		c.secure = true
	}
}

// MaxAge sets the cookie's max age.
func MaxAge(age time.Duration) cookieOption {
	return func(c *cookie) {
		c.maxAge = option.Some(age)
	}
}

// Expires sets the cookie's expiration.
func Expires(expires time.Time) cookieOption {
	return func(c *cookie) {
		c.expires = option.Some(expires)
	}
}

// Domain sets the cookie's domain.
func Domain(value string) cookieOption {
	if value == "" {
		panic("empty domain value")
	}
	return func(c *cookie) {
		c.domain = value
	}
}

// Path sets the cookie's path.
func Path(value string) cookieOption {
	if value == "" {
		panic("empty path value")
	}
	return func(c *cookie) {
		c.path = value
	}
}

// SameSite sets cookie's same-site.
func SameSite(policy SameSitePolicy) cookieOption {
	if policy == "" {
		panic("invalid same-site")
	}
	return func(c *cookie) {
		c.sameSite = policy
	}
}
