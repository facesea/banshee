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
	// Routes
	router := httprouter.New()
	router.GET("/api/config", getConfig)
	router.GET("/api/projects", getProjects)
	router.GET("/api/project/:id", getProject)
	router.POST("/api/project", createProject)
	router.PATCH("/api/project/:id", updateProject)
	router.DELETE("/api/projects/:id", deleteProject)
	router.GET("/api/project/:id/rules", getProjectRules)
	router.GET("/api/project/:id/users", getProjectUsers)
	router.POST("/api/project/:id/user", addProjectUser)
	router.DELETE("/api/project/:id/user/:user_id", deleteProjectUser)
	router.GET("/api/users", getUsers)
	router.GET("/api/user/:id", getUser)
	router.POST("/api/user", createUser)
	router.DELETE("/api/user/:id", deleteUser)
	router.PATCH("/api/user/:id", updateUser)
	router.POST("/api/rule", createRule)
	router.DELETE("/api/rule/:id", deleteRule)
	router.GET("/api/metric/indexes", getMetricIndexes)
	router.GET("/api/metric/data/:name/:start/:stop", getMetrics)
	// Serve
	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Webapp.Port)
	log.Info("serve on %s..", addr)
	http.ListenAndServe(addr, router)
}
