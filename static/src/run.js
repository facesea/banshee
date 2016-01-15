/**
 * Created by Panda on 16/1/13.
 */
/*@ngInject*/
module.exports = function ($rootScope) {

  $rootScope.$on('$stateChangeError',
    function (event, toState, toParams, fromState, fromParams, error) {
      console.error('$stateChangeError,toState:%s', toState.name);
      console.error(error);
      event.preventDefault();
    });

};
