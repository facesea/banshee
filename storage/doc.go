// Copyright 2015 Eleme Inc. All rights reserved.

/*

Package storage implements banshee's persistence storage.

Structure

Storage handles 4 kinds of data: index, metric, state and admin, the directory
structure is:

	storage/
	    |--index/               --- Metric index          LevelDB
	    |--metric/              --- Metric data           LevelDB
	    |--state-288x300/       --- Detection state       LevelDB
	    |--admin                --- Rules/Users/Projects  SQLite3

The storage directory will be created if not exists.

For each child database, see its package's documentation for more.

*/
package storage
