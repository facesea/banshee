// Copyright 2015 Eleme Inc. All rights reserved.

package webapp

import (
	"github.com/eleme/banshee/models"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"github.com/mattn/go-sqlite3"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// getUsers returns all users.
func getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var users []models.User
	if err := db.Admin.DB().Find(&users).Error; err != nil {
		ResponseError(w, NewUnexceptedWebError(err))
		return
	}
	ResponseJSONOK(w, users)
}

// getUser returns user by id.
func getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Params
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		ResponseError(w, ErrUserID)
		return
	}
	// Query db.
	user := &models.User{}
	if err := db.Admin.DB().First(user, id).Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			ResponseError(w, ErrUserNotFound)
			return
		default:
			ResponseError(w, NewUnexceptedWebError(err))
			return
		}
	}
	ResponseJSONOK(w, user)
}

// createUser request
type createUserRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	EnableEmail bool   `json:"enableEmail"`
	Phone       string `json:"phone"`
	EnablePhone bool   `json:"enablePhone"`
	Universal   bool   `json:"universal"`
}

// createUser creats a user.
func createUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Request
	req := &createUserRequest{
		EnableEmail: true,
		EnablePhone: true,
		Universal:   false,
	}
	if err := RequestBind(r, req); err != nil {
		ResponseError(w, ErrBadRequest)
		return
	}
	// Validation
	if len(req.Name) == 0 {
		ResponseError(w, ErrUserName)
		return
	}
	if len(req.Email) == 0 || !strings.Contains(req.Email, "@") {
		ResponseError(w, ErrUserEmail)
		return
	}
	if len(req.Phone) != 10 && len(req.Phone) != 11 {
		ResponseError(w, ErrUserPhone)
		return
	}
	if ok, _ := regexp.MatchString("^\\d{10,11}", req.Phone); !ok {
		ResponseError(w, ErrUserPhone)
		return
	}
	// Save
	user := &models.User{
		Name:        req.Name,
		Email:       req.Email,
		EnableEmail: req.EnableEmail,
		Phone:       req.Phone,
		EnablePhone: req.EnablePhone,
		Universal:   req.Universal,
	}
	if err := db.Admin.DB().Create(user).Error; err != nil {
		// Write errors.
		sqliteErr, ok := err.(sqlite3.Error)
		if ok {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintNotNull:
				ResponseError(w, ErrNotNull)
				return
			case sqlite3.ErrConstraintPrimaryKey:
				ResponseError(w, ErrPrimaryKey)
				return
			case sqlite3.ErrConstraintUnique:
				ResponseError(w, ErrDuplicateUserName)
				return
			}
		}
		// Unexcepted.
		ResponseError(w, NewUnexceptedWebError(err))
		return
	}
	ResponseJSONOK(w, user)
}

// deleteUser deletes a user.
func deleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Params
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		ResponseError(w, ErrUserID)
		return
	}
	// Delete.
	if err := db.Admin.DB().Delete(&models.User{ID: id}).Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			ResponseError(w, ErrUserNotFound)
			return
		default:
			ResponseError(w, NewUnexceptedWebError(err))
			return
		}
	}
}

// updateUser request
type updateUserRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	EnableEmail bool   `json:"enableEmail"`
	Phone       string `json:"phone"`
	EnablePhone bool   `json:"enablePhone"`
	Universal   bool   `json:"universal"`
}

// updateUser updates a user.
func updateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Params
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		ResponseError(w, ErrUserID)
		return
	}
	// Request
	req := &updateUserRequest{}
	if err := RequestBind(r, req); err != nil {
		ResponseError(w, ErrBadRequest)
		return
	}
	// Validation
	if len(req.Name) == 0 {
		ResponseError(w, ErrUserName)
		return
	}
	if len(req.Email) == 0 || !strings.Contains(req.Email, "@") {
		ResponseError(w, ErrUserEmail)
		return
	}
	if len(req.Phone) != 10 && len(req.Phone) != 11 {
		ResponseError(w, ErrUserPhone)
		return
	}
	if ok, _ := regexp.MatchString("^\\d{10,11}", req.Phone); !ok {
		ResponseError(w, ErrUserPhone)
		return
	}
	// Find
	user := &models.User{}
	if err := db.Admin.DB().First(user, id).Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			ResponseError(w, ErrUserNotFound)
			return
		default:
			ResponseError(w, NewUnexceptedWebError(err))
			return
		}
	}
	// Patch
	user.Name = req.Name
	user.Email = req.Email
	user.EnableEmail = req.EnableEmail
	user.Phone = req.Phone
	user.EnablePhone = req.EnablePhone
	user.Universal = req.Universal
	if err := db.Admin.DB().Save(user).Error; err != nil {
		if err == gorm.RecordNotFound {
			// User not found.
			ResponseError(w, ErrUserNotFound)
			return
		}
		// Write errors.
		sqliteErr, ok := err.(sqlite3.Error)
		if ok {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintNotNull:
				ResponseError(w, ErrNotNull)
				return
			case sqlite3.ErrConstraintUnique:
				ResponseError(w, ErrDuplicateUserName)
				return
			}
		}
		// Unexcepted error.
		ResponseError(w, NewUnexceptedWebError(err))
		return
	}
	ResponseJSONOK(w, user)
}

// getUserProjects gets usr projects.
func getUserProjects(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Params
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		ResponseError(w, ErrProjectID)
		return
	}
	// Get User.
	user := &models.User{}
	if err := db.Admin.DB().First(user, id).Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			ResponseError(w, ErrUserNotFound)
			return
		default:
			ResponseError(w, NewUnexceptedWebError(err))
			return
		}
	}
	// Query
	var projs []models.Project
	if user.Universal {
		// Get all projects for universal user.
		if err := db.Admin.DB().Find(&projs).Error; err != nil {
			ResponseError(w, NewUnexceptedWebError(err))
			return
		}
	} else {
		// Get related projects for this user.
		if err := db.Admin.DB().Model(user).Association("Projects").Find(&projs).Error; err != nil {
			ResponseError(w, NewUnexceptedWebError(err))
			return
		}
	}
	ResponseJSONOK(w, projs)
}
