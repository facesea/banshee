// Copyright 2015 Eleme Inc. All rights reserved.

package webapp

import (
	"net/http"
	"strings"
)

// Routes start with `/admin` should be authed.
const authPathPrefix = "/admin"

// staticHandler serves all non-api routes as static files.
type staticHandler struct {
	// Auth
	auth *authHandler
	// FileServer
	fileHandler http.Handler
}

// newStaticHandler creates a staticHandler.
func newStaticHandler(root http.FileSystem, auth *authHandler) *staticHandler {
	fileHandler := http.FileServer(root)
	return &staticHandler{auth, fileHandler}
}

// ServeHTTP implements http.Handler.
func (sh *staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Auth /admin prefixed routes.
	if strings.HasPrefix(r.URL.Path, authPathPrefix) && !sh.auth.auth(w, r) {
		return
	}
	// Public resource.
	sh.fileHandler.ServeHTTP(w, r)
}
