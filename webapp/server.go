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

// Globals
var (
	// Config
	cfg *config.Config
	// Storage
	db *storage.DB
)

// Init globals.
func Init(c *config.Config, d *storage.DB) {
	cfg = c
	db = d
}

// Start http server.
func Start(c *config.Config, d *storage.DB) {
	// Init globals.
	cfg = c
	db = d
	// Auth
	auth := newAuthHandler(cfg.Webapp.Auth[0], cfg.Webapp.Auth[1])
	// Routes
	router := httprouter.New()
	// Api
	router.GET("/api/config", auth.handler(getConfig))
	router.GET("/api/interval", getInterval)
	router.GET("/api/projects", getProjects)
	router.GET("/api/project/:id", getProject)
	router.POST("/api/project", auth.handler(createProject))
	router.PATCH("/api/project/:id", auth.handler(updateProject))
	router.DELETE("/api/project/:id", auth.handler(deleteProject))
	router.GET("/api/project/:id/rules", auth.handler(getProjectRules))
	router.GET("/api/project/:id/users", auth.handler(getProjectUsers))
	router.POST("/api/project/:id/user", auth.handler(addProjectUser))
	router.DELETE("/api/project/:id/user/:user_id", auth.handler(deleteProjectUser))
	router.GET("/api/users", auth.handler(getUsers))
	router.GET("/api/user/:id", auth.handler(getUser))
	router.POST("/api/user", auth.handler(createUser))
	router.DELETE("/api/user/:id", auth.handler(deleteUser))
	router.PATCH("/api/user/:id", auth.handler(updateUser))
	router.GET("/api/user/:id/projects", auth.handler(getUserProjects))
	router.POST("/api/project/:id/rule", auth.handler(createRule))
	router.DELETE("/api/rule/:id", auth.handler(deleteRule))
	router.GET("/api/metric/indexes", getMetricIndexes)
	router.GET("/api/metric/data", getMetrics)
	router.GET("/api/info", getInfo)
	// Static
	router.NotFound = newStaticHandler(http.Dir(cfg.Webapp.Static), auth)
	// Serve
	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Webapp.Port)
	log.Info("webapp is listening and serving on %s..", addr)
	http.ListenAndServe(addr, router)
}
