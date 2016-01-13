// Copyright 2015 Eleme Inc. All rights reserved.

/*

Package admindb handles the admin storage on SQLite3.

Auto Migration

Table schemas will be auto migrated on db opening, that means tables
will be automatically created on the first time banshee starts.

Persistence

Users, Rules and Projects are stored on disk in sqlite3, the
relation between them is:

	User:Project    N:M
	Rule:Project    N:1

To get gorm DB handle:

	adminDBInstance.DB()

Rules Cache

To access rules faster in detector, rules are cached in memory, in a
safemap with RWLock.

To get the rulesCache handle:

	adminDBInstance.RulesCache

*/
package admindb
