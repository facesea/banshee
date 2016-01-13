// Copyright 2015 Eleme Inc. All rights reserved.

/*

Package metricdb handles the storage for metrics.

Design

The db file name is metric and its key-value design in leveldb is:

	Key:   <Name>:<Stamp>
	Value: <Value>:<Score>:<Average>

For less disk usage:

1. Metric stamps will be converted to 36-hex strings before they are put to db.

2. Metric stamps will minus a "stamp horizon" before they are converted to string.

*/
package metricdb
