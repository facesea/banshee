// Copyright 2015 Eleme Inc. All rights reserved.

package models

// Item is rules container which can be a project or what else.
type Item struct {
	Name  string
	Rules []Rule
	Users []User
}
