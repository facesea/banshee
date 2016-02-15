// Copyright 2015 Eleme Inc. All rights reserved.

package webapp

import (
	"github.com/eleme/banshee/version"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// getVersion returns version.
type versionResponse struct {
	Version string `json:"version"`
}

func getVersion(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	v := &versionResponse{Version: version.Version}
	ResponseJSONOK(w, v)
}
