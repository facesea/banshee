/**
 * Created by Panda on 16/1/13.
 */
/*@ngInject*/
module.exports = function ($locationProvider, $mdThemingProvider) {
  $locationProvider.html5Mode(false);

  $mdThemingProvider.theme('default')
    .primaryPalette('blue', {
      'default': '800'
    });
};
