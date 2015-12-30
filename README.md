Banshee
=======

Anomalies detection system for periodic metrics.

Features
--------

* Designed for periodic metrics.
* Easy to get started.
* Automatically anomalies detection.
* Alerting rules admin panels.
* Detection result visualization.
* Built for a large quantity of metrics.

Get Started
-----------

1. Install [go](https://golang.org/).
2. Install [godep](https://github.com/tools/godep).
3. Run `make`.
4. Run `./banshee -c <path/to/config.json>`.

Configuration
-------------

The default configuration is [exampleConfig.json](config/exampleConfig.json).

Statsd Integration
------------------

1. Install `statsd-banshee` on the statsd servers: `npm install statsd-banshee`.
2. Add module `statsd-banshee` to statsd's backends:

   ```js
   {
   , backends: ['statsd-banshee']
   , bansheeHost: 'localhost'
   , bansheePort: 2015
   }
   ```

Alerting Command
----------------

See [alerter/exampleCommand](alerter/exampleCommand).

Net Protocol
------------

Very simple line-based:

```
<Name> <Timestamp> <Value>\n
```

For example:

```bash
$ telnet 0.0.0.0 2015                                                                                                                                                                                            18 ↵ (go1.5.2 node@v5.0.0)
Trying 0.0.0.0...
Connected to 0.0.0.0.
Escape character is '^]'.
counter.foo 1451471948 3.14
```

Algorithms
-----------

1. Anomalies detection: [3-sigma](https://en.wikipedia.org/wiki/68%E2%80%9395%E2%80%9399.7_rule).

   ```py
   if abs(v-avg) > 3*stdev:
       return True  # anomaly
   ```

2. Metric trendings following: [ewma](https://en.wikipedia.org/wiki/Moving_average) and ewms:

   ```py
   avg_old = avg
   avg = avg*(1-f) + v*f
   stdev = sqrt((1-f)*stdev*stdev + f*(v-avg_old)*(v-avg))
   ```

License
-------

MIT Copyright (c) 2015 - 2016 Eleme, Inc.
