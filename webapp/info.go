// Copyright 2015 Eleme Inc. All rights reserved.

package webapp

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type info struct {
	NumMetric int `json:"numMetric"`
}

func getInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	inf := &info{
		NumMetric: db.Index.Len(),
	}
	ResponseJSONOK(w, inf)
}
