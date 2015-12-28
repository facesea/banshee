// Copyright 2015 Eleme Inc. All rights reserved.

package webapp

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// getConfig returns config.
func getConfig(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ResponseJSON(w, cfg)
}
