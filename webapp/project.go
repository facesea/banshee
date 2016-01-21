// Copyright 2015 Eleme Inc. All rights reserved.

package webapp

import (
	"github.com/eleme/banshee/models"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
)

// getProjects returns all projects.
func getProjects(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var projs []models.Project
	if err := db.Admin.DB().Find(&projs).Error; err != nil {
		ResponseError(w, NewUnexceptedWebError(err))
		return
	}
	ResponseJSONOK(w, projs)
}

// getProject returns project by id.
func getProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Params
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		ResponseError(w, ErrProjectID)
		return
	}
	// Query db.
	proj := &models.Project{}
	if err := db.Admin.DB().First(proj, id).Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			ResponseError(w, ErrProjectNotFound)
			return
		default:
			ResponseError(w, NewUnexceptedWebError(err))
			return
		}
	}
	ResponseJSONOK(w, proj)
}

// createProject request
type createProjectRequest struct {
	Name string `json:"name"`
}

// createProject creates a project.
func createProject(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Request
	req := &createProjectRequest{}
	if err := RequestBind(r, req); err != nil {
		ResponseError(w, ErrBadRequest)
		return
	}
	// Validate
	if err := validateProjectName(req.Name); err != nil {
		ResponseError(w, err)
		return
	}
	// Save.
	proj := &models.Project{Name: req.Name}
	if err := db.Admin.DB().Create(proj).Error; err != nil {
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
				ResponseError(w, ErrDuplicateProjectName)
				return
			}
		}
		ResponseError(w, NewUnexceptedWebError(err))
		return
	}
	ResponseJSONOK(w, proj)
}

// updateProject request
type updateProjectRequest struct {
	Name string `json:"name"`
}

// updateProject updates a project.
func updateProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Params
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		ResponseError(w, ErrProjectID)
		return
	}
	// Request
	req := &updateProjectRequest{}
	if err := RequestBind(r, req); err != nil {
		ResponseError(w, ErrBadRequest)
		return
	}
	// Validate
	if err := validateProjectName(req.Name); err != nil {
		ResponseError(w, err)
		return
	}
	// Find
	proj := &models.Project{}
	if err := db.Admin.DB().First(proj, id).Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			ResponseError(w, ErrProjectNotFound)
			return
		default:
			ResponseError(w, NewUnexceptedWebError(err))
			return
		}
	}
	// Patch.
	proj.Name = req.Name
	if err := db.Admin.DB().Save(proj).Error; err != nil {
		if err == gorm.RecordNotFound {
			// Not found.
			ResponseError(w, ErrProjectNotFound)
			return
		}
		// Writer errors.
		sqliteErr, ok := err.(sqlite3.Error)
		if ok {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintNotNull:
				ResponseError(w, ErrNotNull)
				return
			case sqlite3.ErrConstraintUnique:
				ResponseError(w, ErrDuplicateProjectName)
				return
			}
		}
		ResponseError(w, NewUnexceptedWebError(err))
		return
	}
	ResponseJSONOK(w, proj)
}

// deleteProject deletes a project.
func deleteProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Params
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		ResponseError(w, ErrProjectID)
		return
	}
	// Delete.
	if err := db.Admin.DB().Delete(&models.Project{ID: id}).Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			ResponseError(w, ErrProjectNotFound)
			return
		default:
			ResponseError(w, NewUnexceptedWebError(err))
			return
		}
	}
}

// getProjectRules gets project rules.
func getProjectRules(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Params
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		ResponseError(w, ErrProjectID)
		return
	}
	// Query
	var rules []models.Rule
	if err := db.Admin.DB().Model(&models.Project{ID: id}).Related(&rules).Error; err != nil {
		ResponseError(w, NewUnexceptedWebError(err))
		return
	}
	for i := 0; i < len(rules); i++ {
		rules[i].BuildRepr()
	}
	ResponseJSONOK(w, rules)
}

// getProjectUsers gets project users.
func getProjectUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Params
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		ResponseError(w, ErrProjectID)
		return
	}
	// Query
	var users []models.User
	if err := db.Admin.DB().Model(&models.Project{ID: id}).Association("Users").Find(&users).Error; err != nil {
		ResponseError(w, NewUnexceptedWebError(err))
		return
	}
	// Universals
	var univs []models.User
	if err := db.Admin.DB().Where("universal = ?", true).Find(&univs).Error; err != nil {
		ResponseError(w, NewUnexceptedWebError(err))
		return
	}
	users = append(users, univs...)
	ResponseJSONOK(w, users)
}

// addProjectUserRequest is the request of addProjectUser
type addProjectUserRequest struct {
	Name string `json:"name"`
}

// projectUser is the tempory select result for table `project_users`
type projectUser struct {
	UserID    int `sql:"user_id"`
	ProjectID int `sql:"project_id"`
}

// addProjectUser adds a user to a project.
func addProjectUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Params
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		ResponseError(w, ErrProjectID)
		return
	}
	// Request
	req := &addProjectUserRequest{}
	if err := RequestBind(r, req); err != nil {
		ResponseError(w, ErrBadRequest)
		return
	}
	// Find user.
	user := &models.User{}
	if err := db.Admin.DB().Where("name = ?", req.Name).First(user).Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			ResponseError(w, ErrUserNotFound)
			return
		default:
			ResponseError(w, NewUnexceptedWebError(err))
			return
		}
	}
	if user.Universal {
		ResponseError(w, ErrProjectUniversalUser)
		return
	}
	// Find proj
	proj := &models.Project{}
	if err := db.Admin.DB().First(proj, id).Error; err != nil {
		if err == gorm.RecordNotFound {
			ResponseError(w, ErrNotFound)
			return
		}
		ResponseError(w, NewUnexceptedWebError(err))
		return
	}
	// Note: Gorm only insert values to join-table if the primary key not in
	// the join-table. So we select the record at first here.
	if err := db.Admin.DB().Table("project_users").Where("user_id = ? and project_id = ?", user.ID, proj.ID).Find(&projectUser{}).Error; err == nil {
		ResponseError(w, ErrDuplicateProjectUser)
		return
	}
	// Append user.
	if err := db.Admin.DB().Model(proj).Association("Users").Append(user).Error; err != nil {
		if err == gorm.RecordNotFound {
			// User or Project not found.
			ResponseError(w, ErrNotFound)
			return
		}
		// Duplicate primay key.
		sqliteErr, ok := err.(sqlite3.Error)
		if ok {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				ResponseError(w, ErrDuplicateProjectUser)
				return
			case sqlite3.ErrConstraintPrimaryKey:
				ResponseError(w, ErrDuplicateProjectUser)
				return
			}
		}
		// Unexcepted error.
		ResponseError(w, NewUnexceptedWebError(err))
		return
	}
}

// deleteProjectUser deletes a user from a project.
func deleteProjectUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Params
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		ResponseError(w, ErrProjectID)
		return
	}
	userID, err := strconv.Atoi(ps.ByName("user_id"))
	if err != nil {
		ResponseError(w, ErrUserID)
		return
	}
	// Find user.
	user := &models.User{}
	if err := db.Admin.DB().First(user, userID).Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			ResponseError(w, ErrUserNotFound)
			return
		default:
			ResponseError(w, NewUnexceptedWebError(err))
			return
		}
	}
	// Find proj.
	proj := &models.Project{}
	if err := db.Admin.DB().First(proj, id).Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			ResponseError(w, ErrProjectNotFound)
			return
		default:
			ResponseError(w, NewUnexceptedWebError(err))
			return
		}
	}
	// Delete user.
	if err := db.Admin.DB().Model(proj).Association("Users").Delete(user).Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			ResponseError(w, ErrNotFound)
			return
		default:
			ResponseError(w, NewUnexceptedWebError(err))
			return
		}
	}
}
