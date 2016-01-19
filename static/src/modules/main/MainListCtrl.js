/*@ngInject*/
module.exports = function ($scope, $stateParams, Metric, DateTimes) {
  var chart = require('./chart');
  $scope.dateTimes = DateTimes;

  $scope.limitList = [{
    label: 'Limit1',
    val: 1
  }, {
    label: 'Limit 30',
    val: 30
  }, {
    label: 'Limit 50',
    val: 50
  }, {
    label: 'Limit 100',
    val: 100
  }, {
    label: 'Limit 500',
    val: 500
  }, {
    label: 'Limit 1000',
    val: 1000
  }];

  $scope.sortList = [{
    label: 'Trending Up',
    val: 0
  }, {
    label: 'Trending Down',
    val: 1
  }]

  $scope.filter = {
    datetime: DateTimes[0].seconds,
    limit: $scope.limitList[0].val,
    sort: $scope.sortList[0].val
  }

  chart.init({
    selector: '#cubism-wrap',
    serverDelay: 10 * 1000,
    step: 10 * 1000,
    stop: false
  });

  setIntervalAndRunNow(function () {
    var params = {
      limit: $scope.filter.limit,
      sort: $scope.filter.sort,
    };
    if ($stateParams.project) {
      params.project = $stateParams.project;
    } else {
      params.pattern = $stateParams.pattern;
    }

    chart.remove();
    Metric.getMetricIndexes(params).$promise
      .then(function(res) {
        plot(res);
      });
  }, 10 * 60 * 1000);

  /**
   * Plot.
   */
  function plot(data) {
    var name, i, metrics = [];
    for (i = 0; i < data.names.length; i++) {
      name = data.names[i][0];
      // TODO
      // metrics.push(feed(name, self.refreshTitle));
      metrics.push(feed(name, function(){}));
    }
    return chart.plot(metrics);
  };

  /**
   * Feed metric.
   * @param {String} name
   * @param {Function} cb // function(data)
   * @return {Metric}
   */
  self.feed = function(name, cb) {
    return handlers.chart.metric(function(start, stop, step, callback) {
      var values = [], i = 0;
      // cast to timestamp from date
      start = (+start - delay) / 1000;
      stop = (+stop - delay) / 1000;
      step = +step / 1000;
      // parameters to pull data
      var params = {
        name: name,
        type: options.type,
        start: start,
        stop: stop
      };
      // request data and call `callback` with values
      // data schema: {name: {String}, times: {Array}, vals: {Array}}
      handlers.metric.getData(params, function(err, data) {
        if (err)
          return handlers.error.fatal(err);
        // the timestamps from statsd DONT have exactly steps `10`
        while (start < stop) {
          while (start < data.times[i]) {
            start += step;
            values.push(start > data.times[i] ? data.vals[i] : 0);
          }
          values.push(data.vals[i++]);
          start += step;
        }
        callback(null, values);
        if (cb)
          cb(data);
      });
    }, name);
  };

  function setIntervalAndRunNow(fn, ms) {
    fn();
    return setInterval(fn, ms);
  };
  // Replace this with context.graphite and graphite.metric!
  function random(x) {
    var value = 0,
      values = [],
      i = 0,
      last;
    return context.metric(function (start, stop, step, callback) {
      start = +start, stop = +stop;
      if (isNaN(last)) last = start;
      while (last < stop) {
        last += step;
        value = Math.max(-10, Math.min(10, value + .8 * Math.random() - .4 + .2 * Math.cos(i += x * .02)));
        values.push(value);
      }
      callback(null, values = values.slice((start - stop) / step));
    }, x);
  }
};
