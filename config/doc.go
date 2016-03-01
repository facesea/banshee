// Copyright 2015 Eleme Inc. All rights reserved.

/*

Package config handles configuration parser and container.

Example

Configuration is in JSON file, for example:

	{
	  "interval": 10,
	  "period": 86400,
	  "expiration": 604800,
	  "storage": {
	    "path": "storage/"
	  },
	  "detector": {
	    "port": 2015,
	    "trendingFactor": 0.1,
	    "filterOffset": 0.01,
	    "leastCount": 30,
	    "blacklist": ["statsd.*"],
	    "intervalHitLimit": 100,
	    "defaultThresholdMaxs": {"timer.mean_90": 300},
	    "defaultThresholdMins": {},
	    "fillBlankZeros": ["counter.*.exc"]
	  },
	  "webapp": {
	    "port": 2016,
	    "auth": ["user", "pass"],
	    "static": "static/dist",
		"notice": {}
	  },
	  "alerter": {
	    "command": "",
	    "workers": 4,
	    "interval": 1200,
	    "oneDayLimit": 5,
	    "defaultSilentTimeRange": [0, 6]
	  },
	  "cleaner": {
	    "interval": 10800,
	    "threshold": 259200
	  }
	}

To use config file with banshee:
	banshee -c path/to/config.json


Documents

The documents for each configuration item with default values:

	interval                       // All metrics incoming interval (in seconds), default: 10
	period                         // All metrics period (in seconds), default: 86400 (1 day)
	expiration                     // All metrics expiration (in seconds), default: 604800 (7 days)
	storage.path                   // Storage directory path.
	detector.port                  // Detector tcp port to listen.
	detector.trendingFactor        // Detection weighted moving factor, should be a number between 0 and 1, default: 0.1
	detector.filterOffset          // Offset to filter history data, as a percentage to period, default: 0.01
	detector.leastCount            // Least count to start detection. default: 30
	detector.blacklist             // Incoming metrics blacklist, each one should be a wildcard pattern, default: []
	detector.intervalHitLimit      // Limitation for number of filtered metrics for each rule in one interval. default: 100
	detector.defaultThresholdMaxs  // Default thresholdMax for all rules, a wildcard pattern to number map. default: {}
	detector.defaultThresholdMins  // Default thresholdMin for all rules, a wildcard pattern to number map. default: {}
	detector.filterBlankZeros      // Detector will fill metric blanks with zeros if it matches any of these wildcard patterns. default: []
	webapp.port                    // Webapp http port to listen.
	webapp.auth                    // Webapp admin pages basic auth, in form of [user, pass], use empty string ["", ""] to disable auth. default: ["admin", "admin"]
	webapp.static                  // Webapp static files (htmls/css/js) path, default: "static/dist"
	webapp.notice                  // Webapp notice in HTML, default: {}, example: {"docs": "url-to-docs"}
	alerter.command                // Alerter command or script to execute on anomalies found. default: ""
	alerter.workers                // Number of workers to consume command execution jobs. default: 4
	alerter.interval               // Minimal interval (in seconds) between two alerting message for one metric. default: 1200
	alerter.oneDayLimit            // Limitation for number of alerting times for one metric in a day. default: 5
	alerter.defaultSilentTimeRange // Default silent time range if the project silent is disabled. default: [0, 6] (00:00~06:00)
	cleaner.interval               // Time interval to check outdated data to clean. default: 10800 (4 hours)
	cleaner.threshold              // One metric will be cleaned if the age it incoming exceeds this threshold (in seconds). default: 259200 (3 days)

*/
package config
