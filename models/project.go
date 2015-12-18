// Copyright 2015 Eleme Inc. All rights reserved.

package models

// Project is rules group.
// One project can have many rules and many recievers.
type Project struct {
	// Name
	Name string
	// Project has many rules.
	Rules []Rule
	// Project has many receivers.
	Users []User
}
