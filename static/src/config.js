/**
 * Created by Panda on 16/1/13.
 */
/*@ngInject*/
module.exports = function ($locationProvider, $mdThemingProvider, $translateProvider, toastrConfig) {
  $locationProvider.html5Mode(false);

  // Theme
  $mdThemingProvider.theme('default')
    .primaryPalette('blue', {
      'default': '800'
    });

  // Translate
  $translateProvider
    .useStaticFilesLoader({
      prefix: './languages/locale-',
      suffix: '.json'
    });
  // Toastr
  angular.extend(toastrConfig, {
    preventDuplicates: true
  });
};
