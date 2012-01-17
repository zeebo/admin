package admin

import "net/http"

type Authorizer interface {
	Authorize(*http.Request) AuthResponse
}

type AuthResponse struct {
	Passed   bool
	Error    string
	Username string

	//Key must be marshallable by the json package
	Key interface{}
}

type AuthFunc func(*http.Request) AuthResponse

func (a AuthFunc) Authorize(req *http.Request) AuthResponse {
	return a(req)
}
