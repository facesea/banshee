// Copyright 2015 Eleme Inc. All rights reserved.

package models

// User is the alerter message receiver.
// One user can receive many projects.
type User struct {
	// Name
	Name string
	// Email
	Email       string
	EnableEmail bool
	// Phone
	Phone       string
	EnablePhone bool
}
