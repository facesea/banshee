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
      } else if (rejection.status === 403) {
        return $rootScope.$emit(403);
      }

      // 兼容angular-loading-bar
      rejection.message = rejection.data.message;
      return $q.reject(rejection);
    }
  };
};
