// Copyright 2015 Eleme Inc. All rights reserved.

/*

Package indexdb handles the storage for indexes.

Design

The db file name is index, and the key-value design in leveldb is:

	Key:   <Name>
	Value: <Stamp>:<Score>:<Average>

Cache

To access indexes faster in webapp, indexes are cached in memory, in
a safemap with a RWLock.

Read operations are in cache.

Write operations are to persistence and cache.

*/
package indexdb
