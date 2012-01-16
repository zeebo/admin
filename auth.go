package admin

import "net/http"

type Authorizer interface {
	Authorize(*http.Request) AuthResponse
}

type AuthResponse struct {
	Passed   bool
	Username string
}

type AuthFunc func(*http.Request) AuthResponse

func (a AuthFunc) Authorize(req *http.Request) AuthResponse {
	return a(req)
}
