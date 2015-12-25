// Copyright 2015 Eleme Inc. All rights reserved.

package webapp

import (
	"fmt"
	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/storage"
	"github.com/eleme/banshee/util/log"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Handle context
type context struct {
	cfg *config.Config
	db  *storage.DB
}

// Global context.
var ctx = &context{}

// Init context with config and db.
func Init(cfg *config.Config, db *storage.DB) {
	ctx.cfg = cfg
	ctx.db = db
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// Serve http up.
func Serve() {
	// Routes
	r := httprouter.New()
	r.GET("/", index)
	// Listen and serve
	addr := fmt.Sprintf("0.0.0.0:%d", ctx.cfg.Webapp.Port)
	log.Info("serve on http://%s", addr)
	http.ListenAndServe(addr, r)
}
