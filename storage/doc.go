// Copyright 2015 Eleme Inc. All rights reserved.

/*

Package storage implements banshee's persistence storage.

Structure

Storage handles 4 kinds of data: index, metric and admin, the directory
structure is:

	storage/
	    |--index/               --- Metric index          LevelDB
	    |--metric/              --- Metric data           LevelDB
	    |--admin                --- Rules/Users/Projects  SQLite3

The storage directory will be created if not exists.

For each child database, see its package's documentation for more.

*/
package storage
