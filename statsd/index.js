/**
 * Backend for github.com/etsy/statsd.
 *
 * Optional configs:
 *   bansheeHost         banshee host to connect to. [default: '0.0.0.0']
 *   bansheePort         banshee port to connect to. [default: 2015]
 *   bansheeIgnores      wildcard patterns to ignore stats. [default: ['statsd.*']]
 *   bansheeTimerFields  timer data fields to forward. [default: ['mean_90', 'count_ps']]
 *
 * Supported metric types:
 *
 *   counter_rates, timers, gauges
 */

'use strict';

var net = require('net');
var  minimatch = require('minimatch');

var config;
var debug;
var logger;

var makers =  {
  'counter_rates': dataFromCounterRates,
  'timer_data': dataFromTimerData,
  'gauges': dataFromGauges,
};

function dataFromCounterRates(key, val, timeStamp) {
  var name = util.format('counter.%s', key);
  return [[name, timeStamp, val]];
}

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

function dataFromGauges(key, val, timeStamp) {
  var name = util.format('gauge.%s', key);
  return [[name, timeStamp, val]];
}

function Banshee(options) {}

Banshee.prototype.connect = function(cb) {
  var self = this;
  var options = {
    host: config.bansheeHost,
    port: config.bansheePort,
  };
  this.conn = net.connect(options, function() {
    if (debug)
      logger.log('banshee conneted');
    cb && cb();
  }).on('error', function(err) {
    if (debug)
      logger.log('banshee socket error: ' + err.message + ', disconnecting..');
    self.conn.destroy();
    self.conn = undefined;
  });
  return this;
};

Banshee.prototype.match = function(key) {
  var patterns = config.bansheeIgnores || ['statsd.*'];
  for (var i = 0; i < patterns.length; i++)
    if (minimatch(key, patterns[i]))
      return true;
  return false;
};

Banshee.prototype.send = function(buf, cb) {
  var self = this;
  if (!this.conn)
    return this.connect(function() {
      return self.conn.write(buf, cb);
    });
  return this.conn.write(buf, cb);
};

Banshee.prototype.flush = function(timeStamp, data) {
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

  if (list.length > 0) {
    var buf = '';
    for (i = 0; i < list.length; i++)
      buf += list[i].join(' ') + '\n';
      this.send(buf, function() {
        if (debug) {
          var msg = util.format("sent to banshee: %s", JSON.stringify(list[0]));
          if (list.length > 1)
            msg = util.format("%s, (%d more..)", msg, list.length - 1)
          logger.log(msg);
        }
      });
  }
};

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
