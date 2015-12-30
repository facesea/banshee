// Copyright 2015 Eleme Inc. All rights reserved.

package webapp

import (
	"encoding/json"
	"github.com/eleme/banshee/util/log"
	"io"
	"net/http"
	"strconv"
)

// ResponseJSONOK writes ok response.
func ResponseJSONOK(w http.ResponseWriter, v interface{}) error {
	return ResponseJSON(w, http.StatusOK, v)
}

// ResponseJSON encodes value to json and write as response.
func ResponseJSON(w http.ResponseWriter, code int, v interface{}) error {
	// Encode JSON.
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	// Write response.
	s := string(b)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(s)))
	w.WriteHeader(code)
	io.WriteString(w, s)
	log.Debug("%d - %s", code, s)
	return nil
}

// ResponseError writes WebError as response.
func ResponseError(w http.ResponseWriter, err *WebError) error {
	return ResponseJSON(w, err.Code, err)
}

// RequestBind binds request data into value.
func RequestBind(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
