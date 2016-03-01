/*@ngInject*/
module.exports = function ($resource) {
  return $resource('/api/config', {}, {
    getInterval: {
      method: 'GET',
      url: '/api/interval'
    },
    getNotice: {
      method: 'GET',
      url: '/api/notice'
    }
  });
};
