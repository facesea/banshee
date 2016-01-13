// Copyright 2015 Eleme Inc. All rights reserved.

/*

Package statedb handles the states storage on LevelDB.

Design

The db file name is state-NumGridxGirdLen, and the key-value design
in leveldb is:

	Key:   <Name>:<GridNo>
	Value: <Average>:<StdDev>:<Count>

The GridNo is computed by:

	GirdNo = (Stamp % Period) / GirdLen

*/
package statedb
