// Copyright 2015 Eleme Inc. All rights reserved.

/*

Package config handles configuration parser and container.

Example

Configuration is in JSON file, for example:

	{
	  "interval": 10,
	  "period": 86400,
	  "storage": {
	    "path": "storage/"
	  },
	  "detector": {
	    "port": 2015,
	    "factor": 0.05,
		"leastCount": 30,
	    "blacklist": ["statsd.*"],
	    "intervalHitLimit": 100,
	    "defaultTrustLines": {"timer.count_ps.*": 30},
	    "fillBlankZeros": ["counter.*.exc"]
	  },
	  "webapp": {
	    "port": 2016,
	    "auth": ["user", "pass"],
	    "static": "static/dist"
	  },
	  "alerter": {
	    "command": "",
	    "workers": 4,
	    "interval": 1200,
	    "oneDayLimit": 5
	  }
	}

To use config file with banshee:
	banshee -c path/to/config.json


Documents

The documents for each configuration item with default values:

	interval                   // All metrics incoming interval (in seconds), default: 10
	period                     // All metrics period (in seconds), default: 86400 (1 day).
	storage.path               // Storage directory path.
	detector.port              // Detector tcp port to listen.
	detector.factor            // Detection weighted moving factor, should be a number between 0 and 1, default: 0.05
	detector.leastCount        // Least count to start detection. default: 30
	detector.blacklist         // Incoming metrics blacklist, each one should be a wildcard pattern, default: []
	detector.intervalHitLimit  // Limitation for number of filtered metrics for each rule in one interval. default: 100
	detector.defaultTrustLines // Default trustlines for rules, a wildcard pattern to trustline number map. default: {}
	detector.filterBlankZeros  // Detector will fill metric blanks with zeros if it matches any of these wildcard patterns. default: []
	webapp.port                // Webapp http port to listen.
	webapp.auth                // Webapp admin pages basic auth, in form of [user, pass], use empty string ["", ""] to disable auth. default: ["admin", "admin"]
	webapp.static              // Webapp static files (htmls/css/js) path, default: "static/dist"
	alerter.command            // Alerter command or script to execute on anomalies found. default: ""
	alerter.workers            // Number of workers to consume command execution jobs. default: 4
	alerter.interval           // Minimal interval (in seconds) between two alerting message for one metric. default: 1200
	alerter.oneDayLimit        // Limitation for number of alerting times for one metric in a day. default: 5

*/
package config
