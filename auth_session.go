package admin

import (
	"github.com/zeebo/sign"
	"net/http"
)

type AuthSession struct {
	Username string
	Key      interface{}
}

func (a *AuthSession) add(s sign.Signer, w http.ResponseWriter) error {
	data, err := s.Sign(a)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: data,
	})
	return nil
}

func (a *AuthSession) clear(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name: "auth",
	})
}
