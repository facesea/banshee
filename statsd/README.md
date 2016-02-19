Statsd-Banshee
==============

Package statsd-banshee is the statsd backend for [banshee](https://www.npmjs.com/package/statsd-banshee).

https://godoc.org/github.com/eleme/banshee#hdr-Statsd_Integration

Install
-------

```
npm install statsd-banshee
```

Add it to statsd's backends in config.js:

```js
{
, backends: ['statsd-banshee']
, bansheeHost: 'localhost'
, bansheePort: 2015
}
```

Options
-------

```
 bansheeHost         banshee host to connect to. [default: '0.0.0.0']
 bansheePort         banshee port to connect to. [default: 2015]
 bansheeIgnores      wildcard patterns to ignore stats. [default: ['statsd.*']]
 bansheeTimerFields  timer data fields to forward. [default: ['mean_90', 'count_ps']]
```

Support Metric Types
---------------------

* counter `counter.*`
* timer `timer.mean_90.*`, `timer.upper_90.*`, `timer.count_ps.*`
* gauge `gauge.*`

License
-------

MIT (c) Eleme, Inc.
