package admin

import (
	//	"github.com/zeebo/sign"
	"net/http"
	"time"
)

type Signer interface {
	Sign(interface{}) (string, error)
	Unsign(string, interface{}, time.Duration) error
}

type authSession struct {
	Username string
	Key      interface{}
}

func (a *authSession) Add(s Signer, w http.ResponseWriter) error {
	return nil
}

func (a *authSession) Clear(w http.ResponseWriter) error {
	return nil
}
