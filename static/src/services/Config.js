/*@ngInject*/
module.exports = function ($resource) {
  return $resource('/api/config', {}, {
    get: {
      method: 'GET',
      url: '/api/config'
    }
  });
};
