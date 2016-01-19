/**
 * Created by Panda on 16/1/19.
 */
/*@ngInject*/
module.exports = function ($resource) {
  return $resource('/api/metric/:name/:start/:stop', {name: '@name', start: '@start',  stop: '@stop'}, {
    getMetricValues: {
      method: 'GET',
      url: '/api/metric',
      isArray: true
    },
    getMetricIndexes: {
      method: 'GET',
      url: '/api/metric/indexes',
      isArray: true
    }
  });
};
