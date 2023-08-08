package processor

import "io"

type User struct {
	Username string
	Password string
}

type Input struct {
	User   User
	Header map[string][]string
	Method string
	URL    string
	Body   io.ReadCloser
}

type Output struct {
	StatusLine string
	Header     map[string][]string
	Body       io.ReadCloser
}
