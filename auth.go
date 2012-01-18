package admin

import "net/http"

//Authorizer is a type that the admin will use to authorize requests.
type Authorizer interface {
	Authorize(*http.Request) AuthResponse
}

//AuthResponse is the response an Authorizer must return indicating the status
//of the authorization request. If Passed is true, Error must be "". Key is an
//object that will be passed into future contexts that allows you to identify
//which user is logged in. Username is the display username so you don't have
//to query the database every time.
type AuthResponse struct {
	Passed   bool
	Error    string
	Username string

	//Key must be marshallable by the json package
	Key interface{}
}

//AuthFunc is a handy type to turn functions into Authorizers.
//AunthFunc : Authorizer :: http.HandleFunc : http.Handler
type AuthFunc func(*http.Request) AuthResponse

//Authorize implements the Authorizer interface for the AuthFunc type. Calls
//the underlying function and returns that result.
func (a AuthFunc) Authorize(req *http.Request) AuthResponse {
	return a(req)
}
