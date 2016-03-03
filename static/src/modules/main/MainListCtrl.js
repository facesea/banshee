/*@ngInject*/
module.exports = function($scope, $rootScope, $timeout, $stateParams, $translate, Metric, Config, Project) {
  var chart = require('./chart');
  var cubism;
  var initOpt;
  var isInit = false;

  $rootScope.currentMain = true;
  $scope.projectId = $stateParams.project;
  $scope.notice = null;

  $scope.dateTimes = [
    {
      label: 'MAIN_PAST_NOW',
      seconds: 0
    },
    {
      label: 'MAIN_PAST_3HOURS_AGO',
      seconds: 3 * 3600
    },
    {
      label: 'MAIN_PAST_6HOURS_AGO',
      seconds: 6 * 3600
    },
    {
      label: 'MAIN_PAST_1DAY_AGO',
      seconds: 24 * 3600
    },
    {
      label: 'MAIN_PAST_2DAYS_AGO',
      seconds: 48 * 3600
    },
    {
      label: 'MAIN_PAST_3DAYS_AGO',
      seconds: 3 * 24 * 3600
    },
    {
      label: 'MAIN_PAST_4DAYS_AGO',
      seconds: 4 * 24 * 3600
    },
    {
      label: 'MAIN_PAST_5DAYS_AGO',
      seconds: 5 * 24 * 3600
    },
    {
      label: 'MAIN_PAST_6DAYS_AGO',
      seconds: 6 * 24 * 3600
    },
    {
      label: 'MAIN_PAST_7DAYS_AGO',
      seconds: 7 * 24 * 3600
    }
  ];

  $scope.limitList = [{
    label: 'MAIN_LIMIT_1',
    val: 1
  }, {
    label: 'MAIN_LIMIT_30',
    val: 30
  }, {
    label: 'MAIN_LIMIT_50',
    val: 50
  }, {
    label: 'MAIN_LIMIT_100',
    val: 100
  }, {
    label: 'MAIN_LIMIT_500',
    val: 500
  }, {
    label: 'MAIN_LIMIT_1000',
    val: 1000
  }];

  $scope.sortList = [{
    label: 'MAIN_TREND_UP',
    val: 'up'
  }, {
    label: 'MAIN_TREND_DOWN',
    val: 'down'
  }];

  $scope.typeList = [{
    label: 'MAIN_TYPE_VALUE',
    val: 'v'
  }, {
    label: 'MAIN_TYPE_SCORE',
    val: 'm'
  }];

  $scope.autoComplete = {
    searchText: ''
  };

  initOpt = {
    project: $stateParams.project,
    pattern: $stateParams.pattern,
    datetime: $scope.dateTimes[0].seconds,
    limit: $scope.limitList[2].val,
    sort: $scope.sortList[0].val,
    type: $scope.typeList[0].val,
    status: false
  };

  $scope.filter = angular.copy(initOpt);

  $scope.toggleCubism = function() {
    $scope.filter.status = !$scope.filter.status;
    if (!$scope.filter.status) {
      buildCubism();
    } else {
      cubism.stop();
    }
  };

  $scope.restart = function() {
    $scope.filter = angular.copy(initOpt);

    if ($scope.initProject) {
      $scope.project = $scope.initProject;
      $scope.autoComplete.searchText = $scope.project.name;
    } else {
      $scope.project = '';
      $scope.autoComplete.searchText = '';
    }

    $scope.spinner = true;
    $timeout(function() {
      $scope.spinner = false;
    }, 1000);

    buildCubism();
  };

  $scope.searchPattern = function() {
    $scope.filter.project = '';
    $scope.autoComplete.searchText = '';
    buildCubism();
  };

  $scope.searchProject = function(project) {
    $scope.filter.project = project.id;
    $scope.filter.pattern = '';
    $scope.project = project;
    $scope.projectId = project.id;

    buildCubism();
  };


  $scope.$on('$destroy', function() {
    $rootScope.currentMain = false;
  });

  /**
   * watch filter.
   */
  function watchAll() {
    $scope.$watchGroup(['filter.datetime', 'filter.limit', 'filter.sort', 'filter.type'], function() {
      buildCubism();
    });
  }


  function loadData() {
    Project.getAllProjects().$promise
      .then(function(res) {
        var projectId = parseInt($stateParams.project);
        $scope.projects = res;

        if (projectId) {
          $scope.projects.forEach(function(el) {
            if (el.id === projectId) {
              $scope.autoComplete.searchText = el.name;
              $scope.initProject = el;
              $scope.project = el;
            }
          });
        }
      });

    Config.getInterval().$promise
      .then(function(res) {
        $scope.filter.interval = res.interval;

        setIntervalAndRunNow(buildCubism, 10 * 60 * 1000);

        watchAll();
      });

    Config.getNotice().$promise
    .then(function (res) {
      $scope.notice = res;
    });
  }

  function buildCubism() {
    var params = {
      limit: $scope.filter.limit,
      sort: $scope.filter.sort,
    };
    if ($scope.filter.project) {
      params.project = $scope.filter.project;
    } else {
      params.pattern = $scope.filter.pattern;
    }

    chart.remove();

    isInit = false;

    cubism = chart.init({
      selector: '#chart',
      serverDelay: $scope.filter.datetime * 1000,
      step: $scope.filter.interval * 1000,
      stop: false,
      type: $scope.filter.type
    });

    Metric.getMetricIndexes(params).$promise
      .then(function(res) {
        plot(res);
      });
  }

  /**
   * Plot.
   */
  function plot(data) {
    var name, i, metrics = [];
    for (i = 0; i < data.length; i++) {
      name = data[i].name;
      metrics.push(feed(name, data, refreshTitle));
    }

    return chart.plot(metrics);
  }

  function refreshTitle(data) {
    var _titles = d3.selectAll('.title')[0];
    if (isInit) {
      return;
    }
    _titles.forEach(function(el, index) {
      var _el = _titles[index];
      var currentEl = data[index];
      var className = getClassNameByTrend(currentEl.score);
      var str;
      var _box = ['<div class="box"><span>' + $translate.instant('MAIN_METRIC_RULES_TEXT') + '<span class="icon-tr"></span></span><ul>'];

      for (var i = 0; i < currentEl.matchedRules.length; i++) {
        var rule = currentEl.matchedRules[i];
        _box.push('<li><a href="#/admin/project/' + rule.projectID + '">' + rule.pattern + '</a></li>');
      }
      _box.push('</ul></div>');

      str = [
          '<a href="#/main?pattern=' + currentEl.name + '" class="' + className + '">',
          getTextByTrend(currentEl.score),
          currentEl.name,
          '</a>',
          _box.join('')
      ].join('');

      _el.innerHTML = str;
      isInit = true;
    });
  }

  /**
   * Get title class name.
   * @param {Number} trend
   * @return {String}
   */
  function getClassNameByTrend(trend) {
    if (Math.abs(trend) >= 1) {
      return 'anomalous';
    }
    return 'normal';
  }

  /**
   * Scrollbars
   */
  function initScrollbars() {
    $('.chart-box-top').scroll(function() {
      $('.chart-box').scrollLeft($('.chart-box-top').scrollLeft());
    });
    $('.chart-box').scroll(function() {
      $('.chart-box-top').scrollLeft($('.chart-box').scrollLeft());
    });
  }

  /**
   * Feed metric.
   * @param {String} name
   * @param {Function} cb // function(data)
   * @return {Metric}
   */
  function feed(name, data, cb) {
    return chart.metric(function(start, stop, step, callback) {
      var values = [],
        i = 0;
      // cast to timestamp from date
      start = parseInt((+start - $scope.filter.datetime) / 1000);
      stop = parseInt((+stop - $scope.filter.datetime) / 1000);
      step = parseInt(+step / 1000);
      // parameters to pull data
      var params = {
        name: name,
        start: start,
        stop: stop
      };
      // request data and call `callback` with values
      // data schema: {name: {String}, times: {Array}, vals: {Array}}
      Metric.getMetricValues(params, function(data) {
        // the timestamps from statsd DONT have exactly steps `10`
        var len = data.length;
        while (start < stop && i < len) {

          while (start < data[i].stamp) {
            start += step;
            if ($scope.filter.type === 'v') {
              values.push(start > data[i].stamp ? data[i].value : 0);
            } else {
              values.push(start > data[i].stamp ? data[i].score : 0);
            }
          }

          if ($scope.filter.type === 'v') {
            values.push(data[i++].value);
          } else {
            values.push(data[i++].score);
          }
          start += step;
        }
        callback(null, values);

      });
      cb(data);
    }, name);
  }

  /**
   * Get trend text.
   * @param {Number} trend
   * @return {String}
   */
  function getTextByTrend(trend) {
    if (trend > 0) {
      return '↑ ';
    }

    if (trend < 0) {
      return '↓ ';
    }

    return '- ';
  }

  function setIntervalAndRunNow(fn, ms) {
    fn();
    return setInterval(fn, ms);
  }

  loadData();
  initScrollbars();
};
