/*@ngInject*/
module.exports = function ($scope, $rootScope, $timeout, $stateParams, Metric, Config, Project, DateTimes) {
  var chart = require('./chart');
  var cubism;
  var initOpt;
  var isInit = false;
  $rootScope.currentMain = true;
  $scope.dateTimes = DateTimes;
  $scope.projectId = $stateParams.project;

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
  }];

  $scope.typeList = [{
    label: 'Value',
    val: 'v'
  }, {
    label: 'Score',
    val: 'm'
  }];

  $scope.autoComplete = {
    searchText: ''
  };

  initOpt = {
    project: $stateParams.project,
    pattern: $stateParams.pattern,
    datetime: DateTimes[0].seconds,
    limit: $scope.limitList[2].val,
    sort: $scope.sortList[0].val,
    type: $scope.typeList[0].val,
    status: false
  };

  $scope.filter = angular.copy(initOpt);

  $scope.toggleCubism = function () {
    $scope.filter.status = !$scope.filter.status;
    if (!$scope.filter.status) {
      buildCubism();
    } else {
      cubism.stop();
    }
  };

  $scope.restart = function () {
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

  $scope.searchPattern = function () {
    $scope.filter.project = '';
    $scope.autoComplete.searchText = '';
    buildCubism();
  };

  $scope.searchProject = function (project) {
    $scope.filter.project = project.id;
    $scope.filter.pattern = '';
    $scope.project = project;
    $scope.projectId = project.id;

    buildCubism();
  };


  $scope.$on('$destroy', function () {
    $rootScope.currentMain = false;
  });

  function loadData() {
    Project.getAllProjects().$promise
      .then(function (res) {
        var projectId = parseInt($stateParams.project);
        $scope.projects = res;

        if (projectId) {
          $scope.projects.forEach(function (el) {
            if (el.id === projectId) {
              $scope.autoComplete.searchText = el.name;
              $scope.initProject = el;
              $scope.project = el;
              setTitle();
            }
          });
        }
      });

    Config.getInterval().$promise
      .then(function (res) {
        $scope.filter.interval = res.interval;

        setIntervalAndRunNow(buildCubism, 10 * 60 * 1000);

        watchAll();
      });
  }
  /**
   * watch filter.
   */
  function watchAll() {
    $scope.$watchGroup(['filter.datetime', 'filter.limit', 'filter.sort', 'filter.type'], function () {
      buildCubism();
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

    setTitle();

    isInit = false;

    cubism = chart.init({
      selector: '#chart',
      serverDelay: $scope.filter.datetime * 1000,
      step: $scope.filter.interval * 1000,
      stop: false,
      type: $scope.filter.type
    });

    Metric.getMetricIndexes(params).$promise
      .then(function (res) {
        plot(res);
      });
  }

  function setTitle() {
    if ($scope.filter.project && $scope.project) {
      $scope.title = 'Project: ' + $scope.project.name;
      $scope.showLink = true;
      return;
    }

    $scope.showLink = false;

    if ($scope.filter.pattern) {
      $scope.title = 'Pattern: ' + $scope.filter.pattern;
      return;
    }

    $scope.title = 'Pattern: *';
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
    _titles.forEach(function (el, index) {
      var _el = _titles[index];
      var currentEl = data[index];
      var className = getClassNameByTrend(currentEl.score);
      var str = [
        '<a href="#/main?pattern=' + currentEl.name + '" class="' + className + '">',
        getTextByTrend(currentEl.score),
        currentEl.name,
        '</a>'
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
    $('.chart-box-top').scroll(function () {
      $('.chart-box').scrollLeft($('.chart-box-top').scrollLeft());
    });
    $('.chart-box').scroll(function () {
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
    return chart.metric(function (start, stop, step, callback) {
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
      Metric.getMetricValues(params, function (data) {
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
