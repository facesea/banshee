// Copyright 2015 Eleme Inc. All rights reserved.

package webapp

import (
	"github.com/eleme/banshee/health"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func getInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ResponseJSONOK(w, health.Get())
}
