// Copyright 2015 Eleme Inc. All rights reserved.

package webapp

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// authHandler provides basic auth utils.
type authHandler struct {
	user string
	pass string
}

// newAuthHandler creates a authHandler.
func newAuthHandler(user, pass string) *authHandler {
	return &authHandler{user, pass}
}

// handler returns a httprouter handler with basic auth protection.
func (a *authHandler) handler(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if a.auth(w, r) {
			h(w, r, ps)
		}
	}
}

func (a *authHandler) auth(w http.ResponseWriter, r *http.Request) bool {
	user, pass, ok := r.BasicAuth()
	if ok && user == a.user && pass == a.pass {
		// Ok
		return true
	}
	// Fail
	w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	return false
}
