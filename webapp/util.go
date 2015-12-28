// Copyright 2015 Eleme Inc. All rights reserved.

package webapp

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

// ResponseJSON encodes value to json and write as response.
func ResponseJSON(w http.ResponseWriter, v interface{}) error {
	// Encode JSON.
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	// Write response.
	s := string(b)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(s)))
	io.WriteString(w, s)
	return nil
}

// ResponseError writes WebError as response.
func ResponseError(w http.ResponseWriter, err *WebError) {
	// Code
	w.WriteHeader(err.code)
	// Message
	io.WriteString(w, err.Error())
}

// RequestBind binds request data into value.
func RequestBind(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
