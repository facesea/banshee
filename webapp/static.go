// Copyright 2015 Eleme Inc. All rights reserved.

package webapp

import (
	"net/http"
	"strings"
)

const authPathPrefix = "/admin"

type staticHandler struct {
	// Auth
	user string
	pass string
	// FileServer
	fileHandler http.Handler
}

func newStaticHandler(root http.FileSystem) *staticHandler {
	user := cfg.Webapp.Auth[0]
	pass := cfg.Webapp.Auth[1]
	fileHandler := http.FileServer(root)
	return &staticHandler{user, pass, fileHandler}
}

// ServeHTTP implements http.Handler.
func (sh *staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Auth /admin prefixed routes.
	if strings.HasPrefix(r.URL.Path, authPathPrefix) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != sh.user || pass != sh.pass {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	}
	sh.fileHandler.ServeHTTP(w, r)
}
