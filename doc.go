// Copyright 2015 Eleme Inc. All rights reserved.

/*

Banshee is a real-time anomalies(outliers) detection system for periodic
metrics.

Use Case

We are using it to monitor our website and rpc services intefaces, including
called frequency, response time and exception calls. Our services send statistics
to statsd, statsd aggregates them every 10 seconds and broadcasts the results to
its backends including banshee, banshee analyzes current metrics with history
data, calculates the trending and alerts us if the trending behaves anomalous.

For example, we have an api named get_user, this api's response time (in
milliseconds) is reported to banshee from statsd every 10 seconds:
	20, 21, 21, 22, 23, 19, 18, 21, 22, 20, ..., 300

Banshee will catch the latest metric 300 and report it as an anomaly.

Why don't we just set a fixed threshold instead (i.e. 200ms)? This may also works
but it is boring and hard to maintain a lot of thresholds. Banshee will analyze
metric trendings automatically, it will find the "thresholds" automatically.

Features

1. Designed for periodic metrics. Reality metrics are always with periodicity,
banshee only peeks metrics with the same "phase" to detect.

2. Multiple alerting rule configuration options, to alert via fixed-thresholds or
via anomalous trendings.

3. Coming with anomalies visualization webapp and alerting rules admin panels.

4. Require no extra storage services, banshee handles storage on disk by itself.

Requirements

Go >= 1.4 and godep.

Installation

Just clone it and run make:

	$ git clone git://github.com/eleme/banshee.git
	$ cd banshee
	$ go get github.com/tools/godep
	$ make

Or simply download pre-compiled banshee binary from
https://github.com/eleme/banshee/releases.

Command Line

Usage:
	banshee [-c config] [-d] [-v]

Flags:
	-c config
		Load config from file.
	-d
		Turn on debug mode.
	-v
		Show version.

Configuration

See package config.

Statsd Integration

In order to forward metrics to banshee from statsd, we need to add the
npm module statsd-banshee to statsd's banckends:

1. Install statsd-banshee on your statsd servers:

	$ cd path/to/statsd
	$ npm install statsd-banshee

2. Add module statsd-banshee to statsd's backends in config.js:

	{
	, backends: ['statsd-banshee']
	, bansheeHost: 'localhost'
	, bansheePort: 2015
	}

Compontents

Banshee have 4 compontents and they are running in the same process:

1. Detector is to detect incoming metrics with history data and store the
results.

2. Webapp is to visualize the detection results and provides panels to manage
alerting rules, projects and users.

3. Alerter is to send sms and emails once anomalies are found.

4. Cleaner is to clean outdated metrics from storage.

Alerting Sender

See package alerter and alerter/exampleCommand.

Implementation Details

1. Detection algorithms, see package detector.

2. Detector input net protocol, see package detector.

3. Storage, see package storage.

4. Filter, see package filter.

License

MIT (c) eleme, inc.

*/
package main
