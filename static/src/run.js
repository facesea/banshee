/**
 * Created by Panda on 16/1/13.
 */
/*@ngInject*/
module.exports = function ($rootScope, $translate, toastr, Config) {

  $rootScope.$on('$stateChangeError',
    function (event, toState, toParams, fromState, fromParams, error) {
      console.error('$stateChangeError,toState:%s', toState.name);
      console.error(error);
      event.preventDefault();
    });

  $rootScope.$on('api.response.error', function (event, rejection) {
    toastr.error(rejection.data.msg);
  });

  Config.getLanguage().$promise
    .then(function (res) {
      $translate.use(res.language);
    });

};
