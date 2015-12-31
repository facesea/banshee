// Backend for github.com/etsy/statsd.
//
// Optional configs:
//   bansheeHost         banshee host to connect to. [default: '0.0.0.0']
//   bansheePort         banshee port to connect to. [default: 2015]
//   bansheeIgnores      wildcard patterns to ignore stats. [default: ['statsd.*']]
//   bansheeTimerFields  timer data fields to forward. [default: ['mean_90', 'count_ps']]
//
// Supported metric types:
//
//   counter_rates, timers, gauges

'use strict';

var net = require('net');
var util = require('util');

// Match string with pattern.
function match(p, s) {
  var arr = p.split("*");
  var i, j;
  for (i = 0, j = 0; i < arr.length; i++) {
    j = s.indexOf(arr[i], j);
    if (j < 0)
      return false;
    j += arr[i].length;
  }
  return true;
}

// Globals
var config;
var debug;
var logger;

// Metric type to makers mappings.
var makers =  {
  'counter_rates': dataFromCounterRates,
  'timer_data': dataFromTimerData,
  'gauges': dataFromGauges,
};

// Make data from counter_rates.
function dataFromCounterRates(key, val, timeStamp) {
  var name = util.format('counter.%s', key);
  return [[name, timeStamp, val]];
}

// Make data from timers.
function dataFromTimerData(key, dict, timeStamp) {
  var fields = config.bansheeTimerFields || ['mean_90', 'count_ps'];
  var data = [];
  for (var i = 0; i < fields.length; i++) {
      var field = fields[i];
      var name = util.format('timer.%s.%s', field, key);
      data.push([name, timeStamp, dict[field]]);
  }
  return data
}

// Make data from gauges.
function dataFromGauges(key, val, timeStamp) {
  var name = util.format('gauge.%s', key);
  return [[name, timeStamp, val]];
}

// Banshee.
function Banshee() {}

// Connect to banshee server.
Banshee.prototype.connect = function(cb) {
  var self = this;
  var options = {
    host: config.bansheeHost || '0.0.0.0',
    port: config.bansheePort || 2015,
  };
  this.conn = net.connect(options, function() {
    // Ok
    if (debug)
      logger.log('banshee conneted');
    cb && cb();
  }).on('error', function(err) {
    // Error
    if (debug)
      logger.log('banshee socket error: ' + err.message + ', disconnecting..');
    self.conn.destroy();
    self.conn = undefined;
  });
  return this;
};

// Match metric.
Banshee.prototype.match = function(key) {
  var patterns = config.bansheeIgnores || ['statsd.*'];
  for (var i = 0; i < patterns.length; i++)
    if (match(patterns[i], key))
      return true;
  return false;
};

// Send buffer.
Banshee.prototype.send = function(buf, cb) {
  var self = this;
  if (!this.conn)
    // Auto connect
    return this.connect(function() {
      return self.conn.write(buf, cb);
    });
  return this.conn.write(buf, cb);
};

// Flush metrics.
Banshee.prototype.flush = function(timeStamp, data) {
  // Collect data.
  var list = [];
  var types = Object.keys(makers);

  for (var i = 0; i < types.length; i++) {
    var type = types[i];
    var dict = data[type];
    for (var key in dict) {
      if (!this.match(key)) {
        var val = dict[key];
        var maker = makers[type];
        if (maker)
          [].push.apply(list, maker(key, val, timeStamp));
      }
    }
  }
  // Make and send buffer.
  if (list.length > 0) {
    var buf = '';
    for (i = 0; i < list.length; i++)
      buf += list[i].join(' ') + '\n';
      // Send
      this.send(buf, function() {
        if (debug)
          logger.log(util.format("sent to banshee %d metrics", list.length));
      });
  }
};

// Export init.
exports.init = function(uptime, _config, events, _logger) {
  logger = _logger || console;
  debug = _config.debug;
  config = _config || {};
  var banshee = new Banshee();
  events.on('flush', function(timeStamp, data) {
    banshee.flush(timeStamp, data);
  });
  banshee.connect();
  return true;
};
