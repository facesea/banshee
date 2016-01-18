/*@ngInject*/
module.exports = function ($q, $rootScope) {
  return {
    response: function (response) {
      return response;
    },
    responseError: function (rejection) {
      //检查是否token是否有效
      if (rejection.status === 401) {
        return $rootScope.$emit(401);
      }

      if(/api/.test(rejection.config.url) && rejection.config.method === 'GET') {
        $rootScope.$emit('api.response.error',rejection);
      }

      // 兼容angular-loading-bar
      rejection.msg = rejection.data.msg;
      return $q.reject(rejection);
    }
  };
};
